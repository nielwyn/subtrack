# Basic Dockerfile for Go Inventory System
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api ./cmd/api

# Final stage - using scratch for minimal image
FROM scratch

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api .

# Copy .env.example for reference
COPY --from=builder /app/.env.example .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]
