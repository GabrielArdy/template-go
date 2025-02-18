.PHONY: generate build init clean build-all build-linux build-windows build-mac deploy

# Build configuration
BINARY_NAME=go-scratch
BINARY_DIR=build
API_SPEC=apis/api.yml

# Go build configuration
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.BuildTime=$(BUILD_TIME)"

# Create build directory
$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

# Generate API code
generate:
	@echo "Generating API client code..."
	@mkdir -p generated
	oapi-codegen \
		--package generated \
		-generate types,server,spec \
		$(API_SPEC) > generated/api.gen.go

# Initialize modules
init:
	@echo "Initializing modules..."
	go mod tidy

# Build target
build: init $(BINARY_DIR)
	@echo "Building for $(GOOS)/$(GOARCH)..."
	go build $(LDFLAGS) \
        -o $(BINARY_DIR)/main \
        cmd/main.go
	@echo "Copying config files..."
	cp ./config/config-*.yml $(BINARY_DIR)/

# Platform specific builds
build-linux:
	@echo "Building for Linux..."
    GOOS=linux GOARCH=amd64 make build

build-windows:
	@echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 make build

build-mac:
	@echo "Building for MacOS..."
    GOOS=darwin GOARCH=amd64 make build

# Build all platforms
build-all: build-linux build-windows build-mac

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BINARY_DIR)

deploy:
	@echo "Deploying to server..."
	gcloud run deploy go-absensi-backend \
	--source . \
	--platform managed \
	--region asia-southeast2 \
	--allow-unauthenticated