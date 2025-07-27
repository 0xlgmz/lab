# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install curl for health checks
RUN apk add --no-cache curl

# Copy the binary and env file from builder
COPY --from=builder /app/main .
COPY .env .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"]
