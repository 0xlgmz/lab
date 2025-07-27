package proxy

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Proxy struct {
	authServiceURL        string
	businessServiceURL    string
	inventoryServiceURL   string
	transactionServiceURL string
	fileServiceURL        string
	menuServiceURL        string
	orderServiceURL       string
	tableServiceURL       string
	logger                *zap.Logger
}

func New(
	authURL, businessURL, inventoryURL, transactionURL,
	fileURL, menuURL, orderURL, tableURL string,
	logger *zap.Logger,
) *Proxy {
	return &Proxy{
		authServiceURL:        authURL,
		businessServiceURL:    businessURL,
		inventoryServiceURL:   inventoryURL,
		transactionServiceURL: transactionURL,
		fileServiceURL:        fileURL,
		menuServiceURL:        menuURL,
		orderServiceURL:       orderURL,
		tableServiceURL:       tableURL,
		logger:                logger,
	}
}

func (p *Proxy) HandleRequest(c *gin.Context) {
	path := c.Request.URL.Path
	method := c.Request.Method

	// Remove /api/v1 prefix from the path
	servicePath := strings.TrimPrefix(path, "/api/v1")

	// Determine which service to route to
	var targetURL string
	switch {
	case strings.HasPrefix(servicePath, "/auth"):
		targetURL = p.authServiceURL
	case strings.HasPrefix(servicePath, "/business") || strings.HasPrefix(servicePath, "/branches"):
		targetURL = p.businessServiceURL
	case strings.HasPrefix(servicePath, "/products") || strings.HasPrefix(servicePath, "/stock"):
		targetURL = p.inventoryServiceURL
	case strings.HasPrefix(servicePath, "/transactions"):
		targetURL = p.transactionServiceURL
	case strings.HasPrefix(servicePath, "/files"):
		targetURL = p.fileServiceURL
	case strings.HasPrefix(servicePath, "/menu") || strings.HasPrefix(servicePath, "/categories") || strings.HasPrefix(servicePath, "/items"):
		targetURL = p.menuServiceURL
	case strings.HasPrefix(servicePath, "/orders"):
		targetURL = p.orderServiceURL
	case strings.HasPrefix(servicePath, "/tables"):
		targetURL = p.tableServiceURL
	default:
		p.logger.Warn("Service not found", zap.String("path", path))
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	// Create target URL
	target, err := url.Parse(targetURL + servicePath)
	if err != nil {
		p.logger.Error("Invalid target URL", zap.Error(err), zap.String("target", targetURL))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid target URL"})
		return
	}

	// Forward the request
	proxyRequest(c, target, method)
}

func proxyRequest(c *gin.Context, target *url.URL, method string) {
	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
		return
	}

	// Create new request with the body
	req, err := http.NewRequest(method, target.String(), bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward query parameters
	req.URL.RawQuery = c.Request.URL.RawQuery

	// Make the request with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Set status code
	c.Writer.WriteHeader(resp.StatusCode)

	// Copy response body
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
