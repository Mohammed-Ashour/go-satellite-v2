# Variables
BINARY_NAME=go_satellite_v2
VERSION=$(shell cat VERSION 2>/dev/null || echo "0.0.0")
LATEST_TAG=$(shell git fetch --tags && git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null || echo "v0.0.0")
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all build clean test version tag publish

# Default target
all: clean build test

# Build the application
build:
	@if [ "${v}" = "" ]; then \
		echo "Please specify version: make tag v=X.X.X"; \
		exit 1; \
	fi
	@echo "Building ${BINARY_NAME}..."
	go build -o bin/${BINARY_NAME} -ldflags="-X 'main.Version=${v}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}'" main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Show current version
version:
	@echo "Current version: ${VERSION}"
	@echo "Latest git tag: ${LATEST_TAG}"

# Create and push new git tag
tag:
	@if [ "${v}" = "" ]; then \
		echo "Please specify version: make tag v=X.X.X"; \
		exit 1; \
	fi
	@if [ "v${v}" = "${LATEST_TAG}" ]; then \
		echo "Error: Version v${v} already exists as a tag"; \
		echo "Latest tag: ${LATEST_TAG}"; \
		exit 1; \
	fi
	@echo "Creating and pushing tag v${v}..."
	@echo "${v}" > VERSION
	git add VERSION
	git commit -m "Bump version to v${v}"
	git tag -a "v${v}" -m "Release version v${v}"	
	git push origin "v${v}"
	git push origin master

# Publish to pkg.go.dev (this happens automatically when pushing a new tag)
publish: tag
	@echo "Publishing v${v} to pkg.go.dev..."
	GOPROXY=proxy.golang.org go list -m github.com/Mohammed-Ashour/go-satellite-v2@v${v}

# Help target
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  version  - Show current version"
	@echo "  tag v=X.X.X    - Create and push new git tag"
	@echo "  publish v=X.X.X - Publish new version to pkg.go.dev"