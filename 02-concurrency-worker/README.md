# 高并发任务系统 (goroutine + worker pool)

本项目是Go实战入门路线的第二阶段，目标是掌握Go的核心优势——并发模型。

## 项目目标

- Worker Pool（可动态扩容）
- goroutine 任务调度
- channel 队列
- 任务状态跟踪
- HTTP + WebSocket 推送任务进度
- 防止 goroutine 泄漏
- context 控制超时与取消

## 学习重点

- Go 与 Java 线程模型的根本差异
- CSP 并发模型
- 如何写不 panic、不泄漏、不阻塞的 worker 系统

## 技术栈

- **Web框架**: Gin
- **并发模型**: 原生 goroutine 和 channel
- **WebSocket**: Gorilla WebSocket
- **配置管理**: Viper
- **日志**: Zap
- **前端**: Vanilla JavaScript + HTML/CSS

## 项目结构

```
02-concurrency-worker/
├── cmd/
│   └── workerpool/
│       └── main.go              # 应用入口
├── internal/
│   ├── api/                     # API接口
│   │   ├── routes.go            # 路由定义
│   │   ├── handlers.go          # HTTP处理器
│   │   └── dto/                 # 数据传输对象
│   ├── service/                 # 业务逻辑
│   │   ├── task_manager.go      # 任务管理器
│   │   ├── worker_pool.go       # Worker池管理
│   │   ├── worker.go            # Worker实现
│   │   ├── task_registry.go     # 任务注册表
│   │   └── *_test.go           # 单元测试
│   ├── model/                   # 数据模型
│   │   ├── task.go              # 任务模型
│   │   └── task_test.go        # 任务模型测试
│   ├── config/                  # 配置管理
│   │   ├── config.go            # 配置实现
│   │   └── config_test.go      # 配置测试
│   ├── app/                     # 应用初始化
│   │   └── app.go               # 应用主逻辑
│   └── pkg/                     # 公共包
│       ├── logger/              # 日志包（复用第一阶段）
│       └── websocket/           # WebSocket相关
│           ├── hub.go           # WebSocket中心
│           └── hub_test.go     # WebSocket测试
├── configs/
│   └── config.yaml              # 配置文件
├── web/                         # 前端页面
│   ├── index.html               # 任务监控页面
│   └── js/                     # JavaScript代码
├── go.mod                       # Go模块文件
├── go.sum                       # Go依赖校验和
└── README.md                    # 项目说明
```

## 核心特性

### 1. 动态Worker Pool
- 根据任务负载自动扩展和收缩Worker数量
- 防止goroutine泄漏的优雅关闭机制
- 支持配置最小和最大Worker数量

### 2. 任务管理系统
- 任务状态跟踪（pending, running, completed, failed, cancelled）
- 任务进度监控
- 任务取消机制
- 任务结果存储

### 3. 实时通信
- WebSocket实时推送任务状态更新
- HTTP REST API接口
- 前端监控面板

### 4. 并发安全保障
- Context控制任务超时和取消
- Channel实现goroutine间安全通信
- Mutex保护共享资源访问
- Panic恢复机制防止程序崩溃

## 启动方式

### 环境要求
- Go 1.21+ (本地开发)
- Docker & Docker Compose (容器化部署)

### 本地开发

```bash
# 克隆项目
cd 02-concurrency-worker

# 安装依赖
go mod tidy

# 运行项目
go run cmd/workerpool/main.go
```

### Docker容器化部署

```bash
# 构建并启动服务
docker-compose up --build

# 后台运行
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

## API端点

### HTTP REST API

- `GET /health` - 健康检查
- `POST /api/v1/tasks` - 提交新任务
- `GET /api/v1/tasks/{id}` - 获取任务详情
- `GET /api/v1/tasks` - 获取所有任务
- `DELETE /api/v1/tasks/{id}` - 取消任务
- `GET /api/v1/tasks/status/stats` - 获取任务统计

### WebSocket API

- `GET /api/v1/ws/tasks` - WebSocket连接端点

## 任务类型

系统支持以下内置任务类型：

1. **example_task** - 示例任务，快速执行
2. **long_running_task** - 长时间运行任务，模拟复杂计算
3. **default_task** - 默认任务，处理任意负载

## 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  port: 8080                    # 服务器端口
  host: "localhost"             # 服务器主机
  mode: "debug"                 # 运行模式

worker:
  min_workers: 5                # 最小Worker数量
  max_workers: 50               # 最大Worker数量
  enable_auto_scaling: true     # 是否启用自动扩展
  scale_up_threshold: 0.8       # 扩展阈值
  scale_down_threshold: 0.3     # 缩减阈值
  scale_check_interval: "30s"   # 扩缩容检查间隔
  shutdown_timeout: "30s"       # 关闭超时时间

task:
  queue_capacity: 1000          # 任务队列容量
  max_concurrent_tasks: 100     # 最大并发任务数
  default_task_timeout: "5m"    # 默认任务超时时间
  cleanup_completed_tasks_after: "1h"  # 清理已完成任务的时间

logger:
  level: "debug"                # 日志级别
  format: "console"             # 日志格式
  output: "stdout"              # 日志输出
```

## 前端监控面板

访问 `http://localhost:8080` 查看任务监控面板，功能包括：

- 实时任务状态显示
- 任务提交表单
- 任务统计图表
- WebSocket连接状态
- 实时日志显示

## 测试

### 运行单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service
go test ./internal/model
```

## 学习要点

### Go并发模型 vs Java线程模型

1. **轻量级**: Goroutines比Java线程更轻量，可以轻松创建成千上万个
2. **通信**: Go使用CSP模型，通过Channel通信而非共享内存
3. **调度**: Go运行时负责Goroutine调度，无需操作系统线程切换开销

### CSP并发模型

- **Communicating Sequential Processes** (通信顺序进程)
- **Don't communicate by sharing memory; share memory by communicating**
- 使用Channel实现goroutine间的安全通信

### 防止goroutine泄漏

1. 使用Context控制生命周期
2. 正确关闭Channel
3. 使用WaitGroup等待goroutine完成
4. 实现优雅关闭机制

## 性能优化

1. **Worker复用**: 避免频繁创建和销毁goroutine
2. **内存管理**: 合理设置channel缓冲区大小
3. **资源清理**: 及时释放已完成任务的资源
4. **负载均衡**: 合理分配任务到各个Worker

## 故障排除

### 常见问题

1. **端口占用**: 确保8080端口未被其他程序占用
2. **依赖问题**: 运行 `go mod tidy` 更新依赖
3. **权限问题**: 确保有足够的权限访问配置文件和日志目录

### 日志查看

```bash
# 查看应用日志
tail -f logs/worker.log

# 查看Docker日志
docker-compose logs -f
```

## 扩展建议

1. 添加数据库持久化存储任务状态
2. 实现任务优先级队列
3. 添加任务依赖关系管理
4. 实现分布式Worker Pool
5. 添加任务重试机制