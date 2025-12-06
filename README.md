# Go 实战入门路线（面向 Java 后端开发者）

本仓库记录了我从 **实战驱动的方式学习 Go（Golang）** 的完整过程。  
我已经有多年 Java 后端经验，但担心未来工作中会遇到 Go，因此本项目作为个人学习路线 + 代码实践仓库。

本路线完全不从语法入门，而是使用 **3 个实战项目 + 1 个工具项目** 逐步掌握 Go 的工程能力。

---

# 📌 学习目标（1 个月战斗力上线）

通过本项目，我的目标是：

- 掌握 Go 在实际工程中的写法（项目结构、错误处理、依赖管理）
- 学会使用 Gin / Gorm / Zap / Viper 等主流库
- 能编写高并发组件（goroutine / channel / worker pool）
- 能开发微服务（gRPC + Protobuf + Etcd）
- 能编写跨平台 CLI 工具
- 能完成生产级 Go 项目（带配置、日志、Docker）

---

# 🧱 项目结构（按阶段推进）

本仓库包含以下 4 个实战项目，每个项目在单独的子目录中进行开发：
```
├── 01-web-api-template/        # Web API 项目（Gin + Gorm）
├── 02-concurrency-worker/      # 并发任务系统（goroutine + worker pool）
├── 03-grpc-microservices/      # 微服务项目（gRPC + Etcd）
└── 04-cli-toolkit/             # CLI 工具（Cobra）
```

每个项目都有自己的说明文档与启动方式。

---

# 🚀 阶段 1：Web API 实战（Gin + Gorm）  
**目录：`01-web-api-template/`**

这一阶段的目标是掌握 Go 在 Web 后端中的基本工程能力。

## 🔧 功能模块

- 使用 Gin 编写 RESTful API
- JWT 用户登录示例
- Gorm ORM + MySQL/PostgreSQL
- Zap 日志中间件
- Viper 配置管理
- 统一响应结构
- 请求参数校验（binding + validator）
- 优雅退出（context + signal）

## 🗂 推荐项目结构（将在本阶段实现）
```
01-web-api-template/
├── cmd/server/main.go
├── internal/
│   ├── api/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── config/
│   ├── middleware/
│   ├── app/
│   └── pkg/
├── pkg/
│   ├── logger/
│   └── db/
├── configs/config.yaml
└── Makefile
```
---

# 🚀 阶段 2：高并发任务系统（goroutine + worker pool）  
**目录：`02-concurrency-worker/`**

第二阶段目标是掌握 Go 的核心优势——并发模型。

## 🔧 将实现的功能

- Worker Pool（可动态扩容）
- goroutine 任务调度
- channel 队列
- 任务状态跟踪
- HTTP + WebSocket 推送任务进度
- 防止 goroutine 泄漏
- context 控制超时与取消

## 🧠 学习重点

- Go 与 Java 线程模型的根本差异
- CSP 并发模型
- 如何写不 panic、不泄漏、不阻塞的 worker 系统

---

# 🚀 阶段 3：gRPC 微服务项目  
**目录：`03-grpc-microservices/`**

这一阶段贴近生产级微服务场景。

## 🔧 将实现的服务

- Product Service（商品服务）
- Order Service（订单服务）
- 使用 gRPC + Protobuf 通信
- Etcd 作为服务发现
- Viper + 环境变量
- Docker Compose 一键启动

## 🧱 目录结构
```
03-grpc-microservices/
├── proto/
├── product-service/
├── order-service/
└── docker-compose.yaml
```
## 📌 学习重点

- gRPC server/client 的结构
- Protobuf 定义规范
- Etcd 的服务注册与发现
- 多服务的本地开发方式

---

# 🚀 阶段 4：跨平台 CLI 工具（Cobra）  
**目录：`04-cli-toolkit/`**

Go 最擅长写 CLI 工具，这个项目用于掌握 Go 的构建能力。

## 🔧 将完成的功能

- CLI 主命令 + 多子命令
- 自动生成 `--help`
- 支持 YAML 配置
- 统一 logging
- 一键编译跨平台二进制（Linux / Mac / Windows）

## 🧱 示例命令结构
```
mycli
├── deploy
├── build
└── config
```
---

# 📦 环境要求

- Go 1.21+
- Docker & Docker Compose
- Protobuf Compiler (`protoc`)
- Git

---

# ▶️ 如何启动学习

按顺序进入各目录，并根据 README 完成本阶段的开发：

1. `01-web-api-template/`  
2. `02-concurrency-worker/`  
3. `03-grpc-microservices/`  
4. `04-cli-toolkit/`

每完成一个阶段，Go 工程能力提升一个等级。

---

# 📝 未来计划

- 添加单元测试（Go mock）
- 集成 OpenTelemetry（Tracing）
- 增加缓存模块（Redis）
- 性能压测与调优

---

# 📚 参考资料

- Gin: https://github.com/gin-gonic/gin  
- Gorm: https://gorm.io/  
- Cobra: https://github.com/spf13/cobra  
- gRPC-Go: https://github.com/grpc/grpc-go  
- Etcd: https://etcd.io/  

---

# 🤝 说明

本仓库主要面向已有 Java 经验、希望快速掌握 Go 工程能力的人。  
目标不是学习语法，而是直接做出 **生产级项目**。

如需补充某个项目模板，请随时提出。
