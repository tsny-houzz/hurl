APP_NAME = hurl
VERSION = 0.1.0
BUILD_DIR = build
SRC_DIR = .
GO_FILES = $(wildcard $(SRC_DIR)/*.go)
BINARY = $(BUILD_DIR)/$(APP_NAME)-macOS

# Go build settings
GOFLAGS = GOOS=darwin GOARCH=amd64

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build: $(BINARY)

$(BINARY): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	$(GOFLAGS) go build -o $(BINARY) $(SRC_DIR)

# Run target
.PHONY: run
run: $(BINARY)
	@$(BINARY)

# Clean target
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
