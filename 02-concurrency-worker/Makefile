# Makefile for Go Concurrency Worker Pool Project

# Variables
BINARY=worker-pool
MAIN_FILE=cmd/workerpool/main.go
BUILD_DIR=build
SOURCE_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Default target
.PHONY: help
help: ## 显示帮助信息
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## 安装项目依赖
	go mod tidy

.PHONY: build
build: ## 编译项目
	go build -o ${BUILD_DIR}/${BINARY} ${MAIN_FILE}

.PHONY: run
run: ## 运行项目
	go run ${MAIN_FILE}

.PHONY: test
test: ## 运行单元测试
	go test -v ./...

.PHONY: test-cover
test-cover: ## 运行单元测试并生成覆盖率报告
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

.PHONY: clean
clean: ## 清理构建产物
	rm -rf ${BUILD_DIR}
	rm -f coverage.out
	rm -f coverage.html

.PHONY: fmt
fmt: ## 格式化代码
	go fmt ./...

.PHONY: vet
vet: ## 检查代码问题
	go vet ./...

.PHONY: lint
lint: ## 运行代码检查工具
	golint ./...

.PHONY: install-tools
install-tools: ## 安装开发工具
	go install golang.org/x/lint/golint@latest

.PHONY: setup
setup: ## 创建必要的目录
	mkdir -p ${BUILD_DIR}
	mkdir -p logs
	mkdir -p data

.PHONY: docker-build
docker-build: ## 构建Docker镜像
	docker build -t ${BINARY}:latest .

.PHONY: docker-run
docker-run: ## 运行Docker容器
	docker run -p 8080:8080 ${BINARY}:latest

.PHONY: all
all: clean setup deps build ## 清理、设置、安装依赖并构建项目

# Display help by default
.DEFAULT_GOAL := help