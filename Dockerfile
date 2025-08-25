FROM golang:1.25-alpine AS builder

RUN apk add --no-cache musl-dev build-base

WORKDIR /src
COPY . .

# Build args for versioning (for ACT/local/dev override)
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=1 go build -tags netgo \
	-ldflags="-linkmode external -extldflags '-static' \
	-X 'github.com/mgmacri/pool-maintenance-app/internal/version.Version=${VERSION}' \
	-X 'github.com/mgmacri/pool-maintenance-app/internal/version.Commit=${COMMIT}' \
	-X 'github.com/mgmacri/pool-maintenance-app/internal/version.BuildDate=${BUILD_DATE}'" \
	-o /pool-maintenance-api ./cmd/main.go

FROM alpine:3.19

COPY --from=builder /pool-maintenance-api /usr/local/bin/pool-maintenance-api
# Copy Swagger docs for Swagger UI
COPY --from=builder /src/docs /usr/local/bin/docs

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/pool-maintenance-api"]
