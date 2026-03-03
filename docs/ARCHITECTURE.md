# 项目架构说明

## 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                             │
│                    (应用入口和组装)                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 初始化并注入依赖
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      Router Layer                           │
│                     (路由配置层)                             │
│  - 定义 URL 路由规则                                         │
│  - 配置中间件                                                │
│  - 将请求分发到对应的 Handler                                │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP 请求
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                          │
│                   (HTTP 处理层/控制器)                       │
│  - 解析 HTTP 请求（参数、请求体）                            │
│  - 数据验证                                                  │
│  - 调用 Service 层处理业务                                   │
│  - 返回 HTTP 响应（JSON）                                    │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 调用业务方法
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                          │
│                     (业务逻辑层)                             │
│  - 核心业务逻辑（密码加密、数据验证）                         │
│  - 事务管理                                                  │
│  - 协调多个 Repository                                       │
│  - 不直接操作数据库                                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 调用数据访问方法
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                         │
│                     (数据访问层)                             │
│  - 封装数据库操作（CRUD）                                    │
│  - 缓存操作（Redis）                                         │
│  - 数据库查询优化                                            │
│  - 不包含业务逻辑                                            │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 使用连接池
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│                    (基础设施层)                              │
│  - 数据库连接池（PostgreSQL）                                │
│  - Redis 客户端                                              │
│  - 消息队列客户端                                            │
│  - 第三方服务客户端                                          │
└─────────────────────────────────────────────────────────────┘
```

## 数据流向

### 1. 创建用户的完整流程

```
前端发送请求
    ↓
POST /api/v1/users
{
  "username": "张三",
  "email": "zhangsan@example.com",
  "password": "123456"
}
    ↓
Router 路由分发
    ↓
Handler.CreateUser()
  - 解析请求体
  - 验证数据格式（username 必填、email 格式、password 最小长度）
    ↓
Service.CreateUser()
  - 检查邮箱是否已存在（调用 Repository.FindByEmail）
  - 加密密码（bcrypt）
  - 创建用户对象
    ↓
Repository.Create()
  - 插入数据库
  - 返回包含 ID 的用户对象
    ↓
Service 返回用户对象
    ↓
Handler 返回 JSON 响应
{
  "id": 1,
  "username": "张三",
  "email": "zhangsan@example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 2. 获取用户的完整流程

```
前端发送请求
    ↓
GET /api/v1/users/1
    ↓
Router 路由分发
    ↓
Handler.GetUser()
  - 获取路径参数 id
    ↓
Service.GetUserByID()
  - 调用 Repository 查询
    ↓
Repository.FindByID()
  - 先查 Redis 缓存
  - 缓存未命中则查数据库
  - 查到后写入缓存
  - 返回用户对象
    ↓
Service 返回用户对象
    ↓
Handler 返回 JSON 响应
```

## 各层职责详解

### Handler Layer（HTTP 处理层）

**职责：**
- 处理 HTTP 协议相关的事情
- 解析请求参数（路径参数、查询参数、请求体）
- 数据验证（使用 binding 标签）
- 调用 Service 层
- 构造 HTTP 响应
- 处理 HTTP 状态码

**不应该做：**
- ❌ 包含业务逻辑
- ❌ 直接操作数据库
- ❌ 复杂的数据处理

**示例：**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }
    
    // 解析和验证请求
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 调用 Service
    user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 返回响应
    c.JSON(201, user)
}
```

### Service Layer（业务逻辑层）

**职责：**
- 实现核心业务逻辑
- 数据验证（业务规则）
- 事务管理
- 协调多个 Repository
- 处理业务异常

**不应该做：**
- ❌ 处理 HTTP 请求/响应
- ❌ 直接操作数据库（通过 Repository）
- ❌ 依赖 HTTP 框架

**示例：**
```go
func (s *UserService) CreateUser(username, email, password string) (*User, error) {
    // 业务规则：检查邮箱是否已存在
    existingUser, _ := s.userRepo.FindByEmail(email)
    if existingUser != nil {
        return nil, errors.New("邮箱已被注册")
    }
    
    // 业务逻辑：加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // 创建用户
    user := &User{
        Username: username,
        Email:    email,
        Password: string(hashedPassword),
    }
    
    // 调用 Repository 保存
    return s.userRepo.Create(user)
}
```

### Repository Layer（数据访问层）

**职责：**
- 封装所有数据库操作
- 缓存操作
- 数据库查询优化
- 数据持久化

**不应该做：**
- ❌ 包含业务逻辑
- ❌ 处理 HTTP 请求
- ❌ 密码加密等业务操作

**示例：**
```go
func (r *UserRepository) FindByID(id string) (*User, error) {
    // 先查缓存
    cached, err := r.cache.Get(ctx, "user:"+id).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 查数据库
    var user User
    if err := r.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    
    // 写缓存
    userData, _ := json.Marshal(user)
    r.cache.Set(ctx, "user:"+id, userData, time.Hour)
    
    return &user, nil
}
```

## 为什么要分层？

### 1. 单一职责原则

每一层只做一件事：
- Handler 只处理 HTTP
- Service 只处理业务
- Repository 只处理数据

### 2. 易于测试

可以单独测试每一层：

```go
// 测试 Service（mock Repository）
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    user, err := service.CreateUser("张三", "test@example.com", "123456")
    // 断言...
}
```

### 3. 易于维护

修改某一层不影响其他层：
- 数据库从 PostgreSQL 换成 MySQL？只改 Repository
- HTTP 框架从 Gin 换成 Echo？只改 Handler
- 业务规则变更？只改 Service

### 4. 代码复用

多个 Handler 可以共享同一个 Service：

```go
// Web API Handler
func (h *UserHandler) GetUser(c *gin.Context) {
    user, _ := h.userService.GetUserByID(id)
    c.JSON(200, user)
}

// gRPC Handler
func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, _ := h.userService.GetUserByID(req.Id)
    return toProtoUser(user), nil
}
```

## 依赖方向

依赖关系是单向的，从上到下：

```
Handler → Service → Repository → Database
```

**不允许反向依赖：**
- ❌ Repository 不能依赖 Service
- ❌ Service 不能依赖 Handler
- ❌ 下层不能依赖上层

## 添加新功能的步骤

假设要添加"订单"功能：

### 1. 定义 Model

```go
// internal/model/order.go
type Order struct {
    ID        uint
    UserID    uint
    ProductID uint
    Amount    float64
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 2. 创建 Repository

```go
// internal/repository/order_repository.go
type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *Order) (*Order, error) {
    // 实现...
}
```

### 3. 创建 Service

```go
// internal/service/order_service.go
type OrderService struct {
    orderRepo *OrderRepository
    userRepo  *UserRepository  // 可能需要验证用户
}

func NewOrderService(orderRepo *OrderRepository, userRepo *UserRepository) *OrderService {
    return &OrderService{
        orderRepo: orderRepo,
        userRepo:  userRepo,
    }
}

func (s *OrderService) CreateOrder(userID, productID uint, amount float64) (*Order, error) {
    // 业务逻辑...
}
```

### 4. 创建 Handler

```go
// internal/handler/order_handler.go
type OrderHandler struct {
    orderService *OrderService
}

func NewOrderHandler(orderService *OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // HTTP 处理...
}
```

### 5. 在 main.go 中注入

```go
func main() {
    // ...
    
    // 创建 Repository
    orderRepo := repository.NewOrderRepository(database.DB)
    
    // 创建 Service
    orderService := service.NewOrderService(orderRepo, userRepo)
    
    // 创建 Handler
    orderHandler := handler.NewOrderHandler(orderService)
    
    // 设置路由
    app := router.SetupRouter(userHandler, orderHandler)
    
    // ...
}
```

### 6. 添加路由

```go
// internal/router/router.go
func SetupRouter(userHandler *UserHandler, orderHandler *OrderHandler) *gin.Engine {
    // ...
    
    orders := api.Group("/orders")
    {
        orders.POST("", orderHandler.CreateOrder)
        orders.GET("/:id", orderHandler.GetOrder)
    }
    
    // ...
}
```

## 总结

这个架构的核心思想：

1. **分层清晰**：每层职责明确
2. **依赖注入**：从外部传入依赖
3. **单向依赖**：只能上层依赖下层
4. **易于测试**：可以 mock 依赖
5. **易于维护**：修改某层不影响其他层

记住：**好的架构不是为了炫技，而是为了让代码更易理解、更易维护、更易扩展**。
