#!/bin/bash

# API测试脚本

# 测试服务器是否运行
echo "Testing server health..."
curl -X GET http://localhost:8080/health
echo ""

# 提交示例任务
echo "Submitting example task..."
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example_task",
    "payload": {
      "message": "Hello, World!",
      "number": 42
    }
  }'
echo ""

# 提交长时任务
echo "Submitting long running task..."
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "long_running_task",
    "payload": {
      "steps": 10
    }
  }'
echo ""

# 获取所有任务
echo "Getting all tasks..."
curl -X GET http://localhost:8080/api/v1/tasks
echo ""

# 获取任务统计
echo "Getting task statistics..."
curl -X GET http://localhost:8080/api/v1/tasks/status/stats
echo ""

echo "API tests completed!"