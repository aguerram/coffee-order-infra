# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/inventory-service

# Production stage
FROM alpine:3.14

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /tmp/inventory-service .

# Expose the port the app runs on
EXPOSE 9000

ENV PORT=9000

# Run the application
CMD ["./inventory-service"]
