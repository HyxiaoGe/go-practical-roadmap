#!/bin/bash

# 设置Go环境变量
export PATH=$HOME/go-install/go/bin:$PATH

# 验证Go安装
echo "Go version: $(go version)"

# 进入项目目录
cd "$(dirname "$0")/../"

echo "Environment setup complete!"
echo "You can now run:"
echo "  go build ./cmd/server/main.go  # Build the application"
echo "  go run ./cmd/server/main.go    # Run the application"
echo "  go test ./...                  # Run all tests"