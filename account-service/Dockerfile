# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY account-service/go.mod ./
COPY account-service/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY account-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install wget for healthcheck
RUN apk add --no-cache wget

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8001

# Run the application
CMD ["./main"] 