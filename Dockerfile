# Build stage
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/api .

# Copy .env if exists (optional)
COPY .env* ./

# Expose port (adjust if your app uses a different port)
EXPOSE 3000

# Run the application
CMD ["./api"]
