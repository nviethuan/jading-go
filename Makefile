APP:=service
PKG:=./
BIN:=./bin/$(APP)
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILDTIME?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
REGISTRY?=ghcr.io/your-org/your-image
PLATFORMS?=linux/amd64,linux/arm64

.PHONY: all lint test build clean docker push

all: lint test build

lint:
	golangci-lint run ./...

test:
	go test ./... -cover -race

build:
	CGO_ENABLED=1 go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILDTIME)" -o $(BIN) $(PKG)

clean:
	rm -rf bin

docker:
	docker buildx build --platform $(PLATFORMS) \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILDTIME=$(BUILDTIME) \
	  -t $(REGISTRY):$(VERSION) -f Dockerfile .

push:
	docker buildx build --push --platform $(PLATFORMS) \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILDTIME=$(BUILDTIME) \
	  -t $(REGISTRY):$(VERSION) -t $(REGISTRY):latest -f Dockerfile .
