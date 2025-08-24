FROM golang:1.25-alpine AS builder

RUN apk add --no-cache musl-dev build-base

WORKDIR /src
COPY . .

# Build a fully static binary
RUN CGO_ENABLED=1 go build -tags netgo -ldflags="-linkmode external -extldflags '-static'" -o /myapp ./cmd/main.go

FROM alpine:3.19

COPY --from=builder /myapp /usr/local/bin/myapp

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/myapp"]
