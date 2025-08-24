# Start from the official Go image for building
FROM golang:1.25.0 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o pool-maintenance-api ./cmd/main.go

# Start a minimal image for running
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/pool-maintenance-api .

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/pool-maintenance-api"]
