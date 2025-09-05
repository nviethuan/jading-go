#!/bin/bash

set -e  # Exit on any error

# Build configuration
CGO_ENABLED=1
GOOS=linux
GOARCH=amd64

# Get version information
VERSION=${VERSION:-"dev-$(date +%Y%m%d-%H%M%S)"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDTIME=$(date +%Y-%m-%dT%H:%M:%SZ)

echo "Building with version: ${VERSION}, commit: ${COMMIT}"

# Build the application
echo "Building application..."
go build -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILDTIME}" \
  -o main .

echo "Build completed successfully!"
