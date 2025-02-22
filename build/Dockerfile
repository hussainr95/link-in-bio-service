# Stage 1: Build
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go application for Linux (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

# Copy the Linux binary from the builder stage
COPY --from=builder /app/main /app/main

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["/app/main"]
