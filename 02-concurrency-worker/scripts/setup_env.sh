#!/bin/bash

# 设置Go环境变量
export PATH=$HOME/go-install/go/bin:$PATH

# 验证Go安装
echo "Go version: $(go version)"

# 进入项目目录
cd "$(dirname "$0")/../"

# 创建必要的目录
echo "Creating necessary directories..."
mkdir -p logs data configs

# 安装依赖
echo "Installing dependencies..."
go mod tidy

echo "Environment setup complete!"
echo "You can now run:"
echo "  make build     # Build the application"
echo "  make run       # Run the application"
echo "  make test      # Run all tests"