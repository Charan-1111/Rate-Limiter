# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download Go modules (optimize cache by downloading before copying all source code)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (disable CGO for alpine compatibility unless specifically needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o rateLimiter main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/rateLimiter .

# Copy deployment configurations required by main.go
COPY --from=builder /app/deploy ./deploy

# Expose the default server port defined in config.json
EXPOSE 8000

# Execute the application
CMD ["./rateLimiter"]