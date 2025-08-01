FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Use a lightweight image for the final stage
FROM alpine:latest

WORKDIR /app

# Install curl for health checks
RUN apk add --no-cache curl

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the .env file
COPY .env .

# Install curl
RUN apk add --no-cache curl

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
