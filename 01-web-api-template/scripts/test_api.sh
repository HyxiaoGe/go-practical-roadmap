#!/bin/bash

# API测试脚本

echo "Testing API endpoints..."

# 测试健康检查端点
echo "1. Testing /health endpoint..."
curl -X GET http://localhost:8080/health
echo -e "\n"

# 测试注册端点
echo "2. Testing /api/v1/register endpoint..."
curl -X POST http://localhost:8080/api/v1/register
echo -e "\n"

# 测试登录端点
echo "3. Testing /api/v1/login endpoint..."
curl -X POST http://localhost:8080/api/v1/login
echo -e "\n"

# 测试用户信息端点（需要JWT令牌）
echo "4. Testing /api/v1/profile endpoint..."
curl -X GET http://localhost:8080/api/v1/profile
echo -e "\n"

echo "API tests completed."