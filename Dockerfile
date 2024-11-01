# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o db-backup

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install required packages
RUN apk add --no-cache \
    postgresql15-client \
    mysql-client \
    sqlite \
    tzdata \
    && mkdir -p /app/backup

# Copy the binary from builder
COPY --from=builder /app/db-backup .

# Create volume for backup files and sqlite database
VOLUME ["/app/backup"]

# Set environment variables
ENV TZ=Asia/Jakarta

# Set entrypoint
ENTRYPOINT ["/app/db-backup"]
