# Garden Logger Makefile

# Variables
BINARY_NAME=garden-logger
BUILD_DIR=.
INSTALL_DIR=$(HOME)/.local/bin
SYSTEM_INSTALL_DIR=/usr/local/bin

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/garden-logger

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Install for current user
.PHONY: install
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

# Install system-wide (requires sudo)
.PHONY: install-system
install-system: build
	sudo cp $(BINARY_NAME) $(SYSTEM_INSTALL_DIR)/
	@echo "Installed to $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)"

# Uninstall from user directory
.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(INSTALL_DIR)/$(BINARY_NAME)"

# Uninstall from system directory
.PHONY: uninstall-system
uninstall-system:
	sudo rm -f $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(SYSTEM_INSTALL_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Test the build
.PHONY: test
test:
	go test ./internal

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Check for issues
.PHONY: vet
vet:
	go vet ./...

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  install        - Install for current user (~/.local/bin)"
	@echo "  install-system - Install system-wide (/usr/local/bin) [requires sudo]"
	@echo "  uninstall      - Remove from ~/.local/bin"
	@echo "  uninstall-system - Remove from /usr/local/bin [requires sudo]"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  fmt            - Format code"
	@echo "  vet            - Check for issues"
	@echo "  help           - Show this help"