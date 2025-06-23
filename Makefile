.PHONY: all build run format clean

GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

BINARY_NAME = stackfetch
MAIN_PACKAGE = ./cmd/stackfetch
VERSION = $(shell git describe --tags 2>/dev/null || echo "latest")

LDFLAGS = -s -w -X main.version=$(VERSION)

all: format build

format:
	@echo "Formatting all Go files..."
	@$(GOFMT) ./...

run:
	@$(GORUN) $(MAIN_PACKAGE)

build:
	@echo "Building $(BINARY_NAME)..."
	@rm -f $(BINARY_NAME)
	@$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: ./$(BINARY_NAME)"
	@ls -alh ./$(BINARY_NAME)

clean:
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

test:
	@echo "Running tests"
	@go test ./...

tidy:
	@$(GOMOD) tidy
	@echo "Go modules tidied"
