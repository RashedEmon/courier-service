FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working dir
WORKDIR /app

# To utilize docker cache
COPY go.mod go.sum ./

COPY app_config.yaml ./

# Download dependencies
RUN go mod tidy

# Copy application source code
COPY . .

# Build application
RUN go build -o main ./cmd/api/

# Minimal base image for container
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
