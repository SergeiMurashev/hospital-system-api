# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY document-service/go.mod ./
COPY proto ./proto/

# Download dependencies and tidy
RUN go mod download && go mod tidy

# Copy source code
COPY document-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install wget for health checks
RUN apk add --no-cache wget

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 50053

# Run the application
CMD ["./main"] 