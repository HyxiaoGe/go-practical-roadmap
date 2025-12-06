# Web API 模板项目 (Gin + Gorm)

本项目是Go实战入门路线的第一阶段，目标是掌握Go在Web后端开发中的基本工程能力。

## 项目目标

- 使用标准库编写 RESTful API
- JWT 用户登录示例
- Gorm ORM + SQLite/MySQL/PostgreSQL
- Zap 日志中间件
- Viper 配置管理
- 统一响应结构
- 请求参数校验
- 优雅退出（context + signal）

## 推荐项目结构

```
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

## 技术栈

- **Web框架**: 标准库http（后续可升级为Gin）
- **ORM**: Gorm
- **日志**: Zap
- **配置管理**: Viper
- **认证**: JWT
- **数据库**: SQLite/PostgreSQL/MySQL（可通过配置切换）

## 启动方式

### 环境要求
- Go 1.21+

### 环境设置

如果系统中没有安装Go，可以使用以下脚本自动下载和设置：

```bash
# 进入项目目录
cd 01-web-api-template

# 运行环境设置脚本
./scripts/setup_env.sh
```

或者手动安装Go：

```bash
# 下载Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# 解压到本地目录
mkdir -p $HOME/go-install
tar -C $HOME/go-install -xzf go1.21.5.linux-amd64.tar.gz

# 设置环境变量
export PATH=$HOME/go-install/go/bin:$PATH
```

### 构建和运行

```bash
# 克隆项目
git clone <repository-url>
cd 01-web-api-template

# 安装依赖
make deps

# 创建必要目录
make setup

# 运行项目
make run

# 或者直接构建
make build
./build/web-api-template
```

### API端点

- `GET /health` - 健康检查
- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `GET /api/v1/profile` - 获取用户信息（需要JWT认证）

### 配置文件

配置文件位于 `configs/config.yaml`，可以根据需要修改以下配置：

- 服务器端口和主机
- 数据库连接信息（支持SQLite、PostgreSQL、MySQL）
- JWT密钥和过期时间
- 日志级别和输出方式

#### 数据库配置示例

1. **SQLite（默认，适合开发环境）**：
```yaml
database:
  driver: "sqlite"
  dsn: "./data/app.db"
```

2. **PostgreSQL（适合生产环境）**：
```yaml
database:
  driver: "postgres"
  dsn: "host=localhost user=postgres password=postgres dbname=webapi port=5432 sslmode=disable"
```

3. **MySQL（适合生产环境）**：
```yaml
database:
  driver: "mysql"
  dsn: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

## 测试

### 运行单元测试

```bash
make test
```

### API端点测试

项目包含一个简单的API测试脚本：

```bash
# 确保服务器正在运行
make run

# 在另一个终端中运行测试脚本
./scripts/test_api.sh
```

## 项目特点

1. **清晰的项目结构**：遵循Go社区标准项目布局
2. **配置管理**：使用Viper支持YAML配置文件和环境变量
3. **日志系统**：集成Zap高性能日志库
4. **数据库访问**：使用Gorm ORM简化数据库操作
5. **认证授权**：JWT令牌认证机制
6. **优雅关闭**：支持信号处理实现平滑重启
7. **Makefile支持**：简化构建和运行过程

## 后续改进计划

- 集成Gin框架替代标准库http
- 添加请求参数验证
- 实现完整的用户服务
- 添加更多中间件（CORS、请求追踪等）
- 完善测试用例
- 添加Docker支持