APP_NAME=pool-maintenance-api
PKG=github.com/mgmacri/pool-maintenance-app/internal/version
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS=-X '$(PKG).Version=$(VERSION)' -X '$(PKG).Commit=$(COMMIT)' -X '$(PKG).BuildDate=$(BUILD_DATE)'

.PHONY: build run clean info

info:
	@echo "Version:     $(VERSION)"
	@echo "Commit:      $(COMMIT)"
	@echo "Build Date:  $(BUILD_DATE)"

build:
	go build -ldflags="$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/main.go

run: build
	./bin/$(APP_NAME)

clean:
	rm -rf bin
