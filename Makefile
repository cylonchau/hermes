# Build configuration
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILD_DIR=cmd
OUT_DIR ?= target
BIN_DIR := $(OUT_DIR)/bin
BUILDOPTS ?= -v

# Get and define default coredns version
COREDNS_VERSION := $(shell cat COREDNS_VERSION)

# Extract modules
modules := $(wildcard $(GOBUILD_DIR)/*)
SUBDIRS := $(patsubst main.go,hermes,$(notdir $(modules)))

.PHONY: all build coredns modules clean help

all: modules coredns

# Build all modules for the current platform
modules:
	@for dir in $(SUBDIRS); do \
		echo "Building module $$dir..."; \
		chmod +x scripts/build.sh && scripts/build.sh $$dir; \
	done

COREDNS_REPO := https://github.com/coredns/coredns.git
COREDNS_DIR := $(OUT_DIR)/coredns-src

# Build CoreDNS with Hermes plugin
# - Read version from COREDNS_VERSION
# - Clone official CoreDNS repository to temporary directory
# - Checkout correct tag
# - Auto-inject hermes plugin if missing
# - Run go generate and build
coredns:
	@echo "Preparing CoreDNS $(COREDNS_VERSION) with Hermes plugin..."
	@mkdir -p $(BIN_DIR)
	@if [ ! -d "$(COREDNS_DIR)" ]; then \
		echo "Cloning CoreDNS repository from $(COREDNS_REPO)..."; \
		git clone $(COREDNS_REPO) $(COREDNS_DIR); \
	fi
	@cd $(COREDNS_DIR) && \
		git fetch origin $(COREDNS_VERSION) && \
		git checkout $(COREDNS_VERSION) && \
		if ! grep -q "hermes:github.com/cylonchau/hermes/plugin" plugin.cfg; then \
			echo "Injecting hermes plugin to plugin.cfg..."; \
			echo "hermes:github.com/cylonchau/hermes/plugin" >> plugin.cfg; \
		fi && \
		go mod edit -replace github.com/cylonchau/hermes=../.. && \
		go generate coredns.go && \
		go mod tidy && \
		echo "Building CoreDNS..." && \
		CGO_ENABLED=0 go build $(BUILDOPTS) -ldflags="-s -w" -o ../bin/coredns-hermes-$(shell go env GOOS)-$(shell go env GOARCH)
	@echo "Done building coredns."

# Build a specific module
build:
	@if [ -z "$(module)" ]; then \
		echo "No module specified. Usage: make build module=<subdir>"; \
		exit 1; \
	fi
	@chmod +x scripts/build.sh && scripts/build.sh $(module)

clean:
	@echo "Cleaning output directory..."
	@rm -rf $(OUT_DIR)
	@echo "Done."

help:
	@echo "Available commands:"
	@echo "  make all             - Build all modules and CoreDNS for current platform"
	@echo "  make modules         - Build all modules in $(GOBUILD_DIR)"
	@echo "  make coredns         - Build CoreDNS with Hermes plugin (auto-configured)"
	@echo "  make build module=X  - Build a specific module (e.g., make build module=hermes)"
	@echo "  make clean           - Remove $(OUT_DIR) directory"
	@echo "  make help            - Show this help message"
