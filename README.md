# Gin 应用 - 手动依赖注入版本

这是一个使用 Gin 框架构建的 Go Web 应用，采用手动依赖注入方式，适合初学者学习。

## 📁 项目结构

```
.
├── cmd/
│   └── server/
│       └── main.go              # 应用入口，负责初始化和启动
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── handler/                 # HTTP 处理层（Controller）
│   │   └── user_handler.go      # 用户相关接口
│   ├── service/                 # 业务逻辑层
│   │   └── user_service.go      # 用户业务逻辑
│   ├── repository/              # 数据访问层
│   │   └── user_repository.go   # 用户数据库操作
│   ├── model/                   # 数据模型
│   │   └── user.go              # 用户模型
│   ├── middleware/              # 中间件
│   │   ├── cors.go              # 跨域中间件
│   │   └── logger.go            # 日志中间件
│   └── router/                  # 路由配置
│       └── router.go            # 路由设置
├── pkg/
│   ├── database/                # 数据库连接
│   │   └── postgres.go          # PostgreSQL 连接池
│   └── redis/                   # Redis 连接
│       └── redis.go             # Redis 客户端
├── config/
│   └── config.yaml              # 配置文件
├── go.mod                       # Go 模块依赖
├── go.sum                       # 依赖校验文件
├── Makefile                     # 常用命令
└── README_CN.md                 # 项目说明（本文件）
```

## 🏗️ 分层架构

项目采用经典的三层架构，依赖关系从上到下：

```
┌─────────────────────────────────────┐
│         Handler Layer               │  ← HTTP 请求处理
│    (处理请求/响应、参数验证)          │
└─────────────────────────────────────┘
              ↓ 依赖
┌─────────────────────────────────────┐
│         Service Layer               │  ← 业务逻辑处理
│    (业务规则、事务管理)               │
└─────────────────────────────────────┘
              ↓ 依赖
┌─────────────────────────────────────┐
│       Repository Layer              │  ← 数据访问
│    (数据库操作、缓存操作)             │
└─────────────────────────────────────┘
              ↓ 依赖
┌─────────────────────────────────────┐
│    Infrastructure Layer             │  ← 基础设施
│    (数据库连接池、Redis 客户端)       │
└─────────────────────────────────────┘
```

## 🔄 手动依赖注入流程

在 `cmd/server/main.go` 中，按以下顺序初始化所有组件：

```go
// 1. 加载配置
config.Load("config/config.yaml")

// 2. 初始化基础设施层（最底层）
database.InitPostgres()  // 数据库连接池
redis.InitRedis()        // Redis 客户端

// 3. 初始化 Repository 层（注入数据库和 Redis）
userRepo := repository.NewUserRepository(database.DB, redis.Client)

// 4. 初始化 Service 层（注入 Repository）
userService := service.NewUserService(userRepo)

// 5. 初始化 Handler 层（注入 Service）
userHandler := handler.NewUserHandler(userService)

// 6. 设置路由（注入 Handler）
app := router.SetupRouter(userHandler)

// 7. 启动服务器
app.Run(":8080")
```

## 🚀 快速开始

### 1. 安装依赖

```bash
# 安装 Go 依赖
make deps

# 或者手动执行
go mod download
go mod tidy
```

### 2. 配置数据库

编辑 `config/config.yaml`，配置数据库连接信息：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: ginapp
  sslmode: disable
  max_idle_conns: 5
  max_open_conns: 25

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

### 3. 启动应用

```bash
# 使用 Makefile
make run

# 或者直接运行
go run cmd/server/main.go
```

应用将在 `http://localhost:8080` 启动。

## 📡 API 接口

### 健康检查

```bash
GET /api/v1/health
```

响应：
```json
{
  "status": "ok"
}
```

### 创建用户

```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "张三",
  "email": "zhangsan@example.com",
  "password": "123456"
}
```

响应：
```json
{
  "id": 1,
  "username": "张三",
  "email": "zhangsan@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 获取用户信息

```bash
GET /api/v1/users/:id
```

响应：
```json
{
  "id": 1,
  "username": "张三",
  "email": "zhangsan@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## 🔧 常用命令

```bash
make help      # 显示帮助信息
make run       # 运行应用（开发模式）
make build     # 编译应用
make clean     # 清理编译产物
make deps      # 安装/更新依赖
make test      # 运行测试
make fmt       # 格式化代码
make vet       # 代码检查
```

## 📚 核心概念

### 1. 依赖注入（Dependency Injection）

依赖注入是一种设计模式，通过构造函数传递依赖，而不是在内部创建：

```go
// ❌ 不好的做法：在内部创建依赖
type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService() *UserService {
    db := database.Connect()  // 内部创建，难以测试
    repo := repository.NewUserRepository(db)
    return &UserService{userRepo: repo}
}

// ✅ 好的做法：通过参数注入依赖
func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}
```

### 2. 连接池（Connection Pool）

数据库和 Redis 连接都使用连接池：

- **不是每次请求都创建新连接**
- 在应用启动时创建连接池
- 请求时从池中获取连接，用完后归还
- 提高性能，减少资源消耗

```go
// 应用启动时创建一次
database.InitPostgres()  // 创建连接池

// 每次请求时复用连接
db.Query("SELECT * FROM users")  // 自动从池中获取连接
```

### 3. 分层架构的好处

- **单一职责**：每层只做自己的事
- **易于测试**：可以 mock 下层依赖
- **易于维护**：修改某一层不影响其他层
- **代码复用**：多个上层可以共享下层

## 🆚 手动注入 vs Wire

### 手动注入的优点

- ✅ 简单直观，容易理解
- ✅ 不需要额外工具
- ✅ 调试方便
- ✅ 适合小型项目

### 手动注入的缺点

- ❌ 依赖多时代码冗长
- ❌ 容易遗漏依赖
- ❌ 修改依赖关系需要手动更新多处

### 何时使用 Wire？

- 项目依赖关系复杂（10+ 层级）
- 团队规模较大
- 微服务架构
- 需要编译时检查依赖

## 📖 学习路径

1. **理解项目结构**：从 `main.go` 开始，看依赖是如何一层层注入的
2. **阅读代码注释**：每个文件都有详细的中文注释
3. **运行项目**：启动应用，测试 API 接口
4. **添加新功能**：尝试添加新的 Handler、Service、Repository
5. **理解分层架构**：体会每一层的职责和好处

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
