# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# 应用名称
BINARY_NAME=web-api-template
MAIN_FILE=./cmd/server/main.go

# 构建目录
BUILD_DIR=build
DATA_DIR=data
LOGS_DIR=logs

# 默认目标
.PHONY: all
all: build

# 构建应用
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed successfully!"

# 运行应用
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	mkdir -p $(DATA_DIR) $(LOGS_DIR)
	$(GOCMD) run $(MAIN_FILE)

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	@echo "Dependencies installed!"

# 清理构建文件
.PHONY: clean
clean:
	@echo "Cleaning build files..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean completed!"

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# 检查代码错误
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# 创建必要的目录
.PHONY: setup
setup:
	@echo "Setting up project directories..."
	mkdir -p $(DATA_DIR) $(LOGS_DIR) $(BUILD_DIR)
	@echo "Setup completed!"

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all     - Build the application (default)"
	@echo "  build   - Build the application"
	@echo "  run     - Run the application"
	@echo "  deps    - Install dependencies"
	@echo "  clean   - Clean build files"
	@echo "  test    - Run tests"
	@echo "  fmt     - Format code"
	@echo "  vet     - Vet code"
	@echo "  setup   - Create necessary directories"
	@echo "  help    - Show this help"

# 默认帮助
.DEFAULT_GOAL := help