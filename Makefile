BINARY_NAME := daedalus
MODULE      := $(shell go list -m 2>/dev/null || echo "unknown")
VERSION     := $(shell node -p "require('./package.json').version" 2>/dev/null || echo "0.0.0")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
DIST_DIR    := dist

LDFLAGS := -s -w \
	-X 'github.com/ricomonster/daedalus/cmd.Version=$(VERSION)' \
	-X 'github.com/ricomonster/daedalus/cmd.Commit=$(COMMIT)' \
	-X 'github.com/ricomonster/daedalus/cmd.BuildTime=$(BUILD_TIME)'

# Detect current OS and ARCH for the install target
OS   := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)

ifeq ($(ARCH),x86_64)
	ARCH := amd64
endif
ifeq ($(ARCH),aarch64)
	ARCH := arm64
endif

INSTALL_SRC := $(DIST_DIR)/$(strip $(BINARY_NAME))-$(strip $(OS))-$(strip $(ARCH))

# ============================================================================
# Default
# ============================================================================

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "  $(BINARY_NAME) — build & install targets"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ \
		{ printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

# ============================================================================
# Dependencies
# ============================================================================

.PHONY: tidy
tidy: ## Run go mod tidy to sync dependencies
	go mod tidy

# ============================================================================
# Build targets
# ============================================================================

.PHONY: build
build: tidy ## Build for the current machine (output: ./release-cli)
	@echo "→ Building for current platform ($(OS)/$(ARCH))..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "✓ Built: ./$(BINARY_NAME)"

.PHONY: build-all
build-all: tidy $(DIST_DIR) ## Build binaries for all platforms into ./dist/
	@echo "→ Building all platform binaries..."
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64  .
	@echo "  ✓ $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64"
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64  .
	@echo "  ✓ $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64"
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64   .
	@echo "  ✓ $(DIST_DIR)/$(BINARY_NAME)-linux-amd64"
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64   .
	@echo "  ✓ $(DIST_DIR)/$(BINARY_NAME)-linux-arm64"
	@echo ""
	@echo "✓ All binaries written to ./$(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/
	@echo "  ✓ $(DIST_DIR)/$(BINARY_NAME)-$(SUFFIX)"

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# ============================================================================
# Install target
# ============================================================================

INSTALL_DIR ?= $(HOME)/.local/bin

.PHONY: install
install: build-all ## Install the correct binary for this machine to ~/.local/bin
	@echo ""
	@echo "→ Detected platform : $(OS)/$(ARCH)"
	@echo "→ Source binary     : $(INSTALL_SRC)"
	@echo "→ Install target    : $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@if [ ! -f "$(INSTALL_SRC)" ]; then \
		echo "✗ Binary not found: $(INSTALL_SRC)"; \
		echo "  Run 'make build-all' first."; \
		exit 1; \
	fi
	@mkdir -p $(INSTALL_DIR)
	@install -m 0755 $(INSTALL_SRC) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@if echo ":$$PATH:" | grep -q ":$(INSTALL_DIR):"; then \
		$(INSTALL_DIR)/$(BINARY_NAME) --version; \
	else \
		echo "⚠  $(INSTALL_DIR) is not in your PATH."; \
		echo "   Add this line to your ~/.zshrc or ~/.bashrc:"; \
		echo ""; \
		echo "     export PATH=\"$(INSTALL_DIR):\$$PATH\""; \
		echo ""; \
		echo "   Then reload your shell:"; \
		echo "     source ~/.zshrc   (or ~/.bashrc)"; \
	fi

.PHONY: uninstall
uninstall: ## Remove the installed binary from ~/.local/bin (or custom INSTALL_DIR)
	@if [ -f "$(INSTALL_DIR)/$(BINARY_NAME)" ]; then \
		rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
		echo "✓ Uninstalled: $(INSTALL_DIR)/$(BINARY_NAME)"; \
	else \
		echo "  Nothing to uninstall — $(INSTALL_DIR)/$(BINARY_NAME) not found."; \
	fi

# ============================================================================
# Release (build-all + zip each binary)
# ============================================================================

.PHONY: release
release: build-all ## Build all binaries and zip each one for distribution
	@echo "→ Packaging binaries..."
	@for bin in $(DIST_DIR)/$(BINARY_NAME)-*; do \
		[ -f "$$bin" ] || continue; \
		zipname="$$bin.zip"; \
		zip -j "$$zipname" "$$bin"; \
		echo "  ✓ $$zipname"; \
	done
	@echo ""
	@echo "✓ Release zips ready in ./$(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/*.zip

# ============================================================================
# Utility
# ============================================================================

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(DIST_DIR) $(BINARY_NAME)
	@echo "✓ Cleaned"

.PHONY: version
version: ## Print the version that will be baked into the binary
	@echo "Version   : $(VERSION)"
	@echo "Commit    : $(COMMIT)"
	@echo "BuildTime : $(BUILD_TIME)"
