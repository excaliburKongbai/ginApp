# 依赖注入详解

## 什么是依赖注入？

依赖注入（Dependency Injection，简称 DI）是一种设计模式，它的核心思想是：

> **不要在对象内部创建依赖，而是从外部传入依赖**

## 为什么需要依赖注入？

### 问题场景

假设我们有一个 `UserService`，它需要使用 `UserRepository` 来访问数据库：

```go
// ❌ 不好的做法：在内部创建依赖
type UserService struct {
    userRepo *UserRepository
}

func NewUserService() *UserService {
    // 在内部创建 Repository
    db := database.Connect()
    repo := NewUserRepository(db)
    
    return &UserService{
        userRepo: repo,
    }
}
```

这种做法有什么问题？

1. **难以测试**：无法替换 `UserRepository` 为 mock 对象
2. **紧耦合**：`UserService` 和 `UserRepository` 的具体实现绑定在一起
3. **难以复用**：如果想用不同的 Repository 实现，必须修改代码
4. **难以维护**：依赖关系隐藏在代码内部，不清晰

### 解决方案：依赖注入

```go
// ✅ 好的做法：通过参数注入依赖
type UserService struct {
    userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}
```

这样做的好处：

1. **易于测试**：可以传入 mock 对象
2. **松耦合**：只依赖接口，不依赖具体实现
3. **易于复用**：可以传入不同的实现
4. **依赖清晰**：从函数签名就能看出依赖关系

## 本项目的依赖注入流程

### 1. 定义构造函数

每一层都定义一个构造函数，接收依赖作为参数：

```go
// Repository 层：接收数据库连接
func NewUserRepository(db *gorm.DB, cache *redis.Client) *UserRepository {
    return &UserRepository{
        db:    db,
        cache: cache,
    }
}

// Service 层：接收 Repository
func NewUserService(userRepo *UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

// Handler 层：接收 Service
func NewUserHandler(userService *UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}
```

### 2. 在 main.go 中组装

在应用启动时，按顺序创建所有对象并注入依赖：

```go
func main() {
    // 1. 初始化基础设施（最底层）
    database.InitPostgres()
    redis.InitRedis()
    
    // 2. 创建 Repository（注入数据库和 Redis）
    userRepo := repository.NewUserRepository(database.DB, redis.Client)
    
    // 3. 创建 Service（注入 Repository）
    userService := service.NewUserService(userRepo)
    
    // 4. 创建 Handler（注入 Service）
    userHandler := handler.NewUserHandler(userService)
    
    // 5. 设置路由（注入 Handler）
    app := router.SetupRouter(userHandler)
    
    // 6. 启动服务器
    app.Run(":8080")
}
```

## 依赖注入的三种方式

### 1. 构造函数注入（推荐）✅

通过构造函数传递依赖：

```go
type UserService struct {
    userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}
```

优点：
- 依赖关系清晰
- 强制提供依赖
- 对象创建后依赖不可变

### 2. 属性注入（不推荐）❌

直接设置对象的属性：

```go
service := &UserService{}
service.userRepo = userRepo  // 直接设置
```

缺点：
- 容易忘记设置依赖
- 依赖可以被修改
- 不清晰

### 3. 方法注入（特殊场景）

通过方法传递依赖：

```go
func (s *UserService) SetRepository(repo *UserRepository) {
    s.userRepo = repo
}
```

使用场景：
- 可选依赖
- 需要动态切换依赖

## 接口与依赖注入

为了更好的解耦，建议使用接口：

```go
// 定义接口
type UserRepository interface {
    FindByID(id string) (*User, error)
    Create(user *User) (*User, error)
}

// Service 依赖接口，不依赖具体实现
type UserService struct {
    userRepo UserRepository  // 接口类型
}

func NewUserService(userRepo UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}
```

好处：
- 可以轻松替换实现
- 方便编写测试（mock）
- 符合依赖倒置原则

## 测试示例

使用依赖注入后，测试变得非常简单：

```go
// 定义 mock Repository
type MockUserRepository struct {
    users map[string]*User
}

func (m *MockUserRepository) FindByID(id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("用户不存在")
    }
    return user, nil
}

// 测试 Service
func TestUserService_GetUserByID(t *testing.T) {
    // 创建 mock Repository
    mockRepo := &MockUserRepository{
        users: map[string]*User{
            "1": {ID: 1, Username: "张三"},
        },
    }
    
    // 注入 mock Repository
    service := NewUserService(mockRepo)
    
    // 测试
    user, err := service.GetUserByID("1")
    if err != nil {
        t.Errorf("期望成功，但返回错误: %v", err)
    }
    if user.Username != "张三" {
        t.Errorf("期望用户名为'张三'，实际为'%s'", user.Username)
    }
}
```

## 常见问题

### Q1: 为什么不直接使用全局变量？

```go
// ❌ 使用全局变量
var DB *gorm.DB

func GetUser(id string) (*User, error) {
    return DB.First(&User{}, id)
}
```

问题：
- 难以测试（无法替换 DB）
- 隐式依赖（不知道函数依赖什么）
- 并发问题（全局变量可能被修改）

### Q2: 依赖太多怎么办？

如果一个对象依赖很多其他对象：

```go
func NewOrderService(
    orderRepo *OrderRepository,
    userRepo *UserRepository,
    productRepo *ProductRepository,
    paymentService *PaymentService,
    notificationService *NotificationService,
) *OrderService {
    // ...
}
```

解决方案：
1. **检查设计**：是否违反单一职责原则？
2. **使用配置对象**：将多个依赖打包
3. **考虑使用 Wire**：自动生成依赖注入代码

### Q3: 什么时候使用 Wire？

手动注入适合：
- 小型项目（< 10 个依赖）
- 学习阶段
- 依赖关系简单

Wire 适合：
- 大型项目（> 10 个依赖）
- 依赖关系复杂
- 团队协作
- 需要编译时检查

## 总结

依赖注入的核心原则：

1. **依赖从外部传入，不在内部创建**
2. **使用构造函数注入**
3. **依赖接口，不依赖实现**
4. **在 main.go 中统一组装**

记住：**依赖注入不是为了炫技，而是为了让代码更易测试、更易维护、更易扩展**。
