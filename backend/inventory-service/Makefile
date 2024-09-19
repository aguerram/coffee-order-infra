# Define a default version if not provided
VERSION ?= 1.0.1

# Build the application locally
build:
	@echo "Building the application for Linux (Alpine)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o tmp/inventory-service .

# Build Docker image using the local build
build-docker: build
	@echo "Building Docker image..."
	@docker build -t coffee-order-demo/inventory-service:$(VERSION) -f Dockerfile.prod .

push-tag:
	@docker tag coffee-order-demo/inventory-service:$(VERSION) demo.goharbor.io/coffee-order-demo/inventory-service:$(VERSION)

push-docker:
	@docker push demo.goharbor.io/coffee-order-demo/inventory-service:$(VERSION)

docker-all: build-docker push-tag push-docker

.PHONY: build build-docker push-tag push-docker docker-all