package cache

import (
	"context"
	"fmt"
	"log"

	"ginApp/internal/config"
	"github.com/redis/go-redis/v9"
)

// Client 全局 Redis 客户端实例
//
// 为什么使用全局变量？
//  1. Redis 客户端内部维护了连接池，应该在整个应用生命周期内只创建一次
//  2. 客户端是线程安全的，可以在多个 goroutine 中并发使用
//  3. 避免重复创建客户端，提高性能
//
// Redis 常见使用场景：
//   - 缓存：减少数据库查询压力，提高响应速度
//   - 会话存储：存储用户登录状态（Session、JWT Token）
//   - 分布式锁：在分布式系统中实现互斥访问
//   - 消息队列：实现异步任务处理
//   - 计数器：实现点赞数、浏览量等实时统计
//   - 排行榜：使用 Sorted Set 实现排名功能
var Client *redis.Client

// InitRedis 初始化 Redis 连接池
//
// 这个函数在应用启动时调用一次（在 main.go 中）
// 它会：
//  1. 读取配置文件中的 Redis 配置
//  2. 创建 Redis 客户端（内部包含连接池）
//  3. 测试连接是否正常
//
// 返回:
//   - error: 如果连接失败返回错误，成功返回 nil
func InitRedis() error {
	// ==================== 第一步：读取配置 ====================
	// 从配置文件中读取 Redis 配置
	// config.Cfg 是全局配置对象，在 main.go 中通过 config.Load() 加载
	cfg := config.Cfg.Redis

	// ==================== 第二步：创建 Redis 客户端 ====================
	// redis.NewClient 创建一个 Redis 客户端
	// 客户端内部维护了连接池，会自动管理连接的创建、复用和释放
	Client = redis.NewClient(&redis.Options{
		// Addr: Redis 服务器地址
		// 格式: "主机:端口"，例如 "localhost:6379"
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),

		// Password: Redis 密码
		// 如果 Redis 没有设置密码，保持为空字符串 ""
		// 生产环境强烈建议设置密码
		Password: cfg.Password,

		// DB: Redis 数据库编号
		// Redis 默认有 16 个数据库（0-15）
		// 不同的数据库之间数据是隔离的
		// 建议：
		//   - 0: 缓存数据
		//   - 1: 会话数据
		//   - 2: 消息队列
		DB: cfg.DB,

		// PoolSize: 连接池大小
		// 表示最多可以创建多少个连接
		// 建议值：10-100（根据实际负载调整）
		// 注意：不要设置太大，避免耗尽 Redis 服务器资源
		PoolSize: cfg.PoolSize,

		// MinIdleConns: 最小空闲连接数
		// 连接池会保持至少这么多空闲连接
		// 好处：减少连接创建的延迟，提高性能
		// 建议值：PoolSize 的 20-50%
		// MinIdleConns: 5,

		// MaxRetries: 最大重试次数
		// 当命令执行失败时，自动重试的次数
		// 建议值：3
		// MaxRetries: 3,

		// DialTimeout: 连接超时时间
		// 建立连接的最大等待时间
		// 建议值：5 秒
		// DialTimeout: 5 * time.Second,

		// ReadTimeout: 读取超时时间
		// 执行命令后等待响应的最大时间
		// 建议值：3 秒
		// ReadTimeout: 3 * time.Second,

		// WriteTimeout: 写入超时时间
		// 发送命令的最大等待时间
		// 建议值：3 秒
		// WriteTimeout: 3 * time.Second,
	})

	// ==================== 第三步：测试连接 ====================
	// context.Background() 创建一个空的上下文
	// Go 的 context 用于控制请求的生命周期、超时、取消等
	// Background() 返回一个永不过期、永不取消的上下文，适合初始化场景
	ctx := context.Background()

	// Ping() 测试 Redis 连接是否正常
	// 它会向 Redis 服务器发送 PING 命令
	// 如果 Redis 正常运行，会返回 PONG 响应
	// 类似于网络中的 "ping" 命令，用于测试连通性
	if err := Client.Ping(ctx).Err(); err != nil {
		// 连接失败，返回错误
		// 可能的原因：
		//   1. Redis 服务器未启动
		//   2. 地址或端口配置错误
		//   3. 密码错误
		//   4. 网络不通
		return fmt.Errorf("连接 Redis 失败: %w", err)
	}

	log.Println("Redis 连接成功")
	return nil
}

// 常用的 Redis 操作示例：
//
// 1. 字符串操作（String）
//   - Set: Client.Set(ctx, "key", "value", time.Hour)  // 设置键值，1小时过期
//   - Get: Client.Get(ctx, "key").Result()             // 获取值
//   - Del: Client.Del(ctx, "key")                      // 删除键
//
// 2. 哈希操作（Hash）
//   - HSet: Client.HSet(ctx, "user:1", "name", "张三") // 设置哈希字段
//   - HGet: Client.HGet(ctx, "user:1", "name")         // 获取哈希字段
//   - HGetAll: Client.HGetAll(ctx, "user:1")           // 获取所有字段
//
// 3. 列表操作（List）
//   - LPush: Client.LPush(ctx, "queue", "task1")       // 从左侧插入
//   - RPop: Client.RPop(ctx, "queue")                  // 从右侧弹出
//
// 4. 集合操作（Set）
//   - SAdd: Client.SAdd(ctx, "tags", "go", "redis")    // 添加成员
//   - SMembers: Client.SMembers(ctx, "tags")           // 获取所有成员
//
// 5. 有序集合操作（Sorted Set）
//   - ZAdd: Client.ZAdd(ctx, "rank", redis.Z{Score: 100, Member: "user1"}) // 添加成员
//   - ZRange: Client.ZRange(ctx, "rank", 0, 9)         // 获取排名前10
//
// 未来可以添加的功能：
// - Close() error                    // 关闭 Redis 客户端（应用退出时调用）
// - Ping() error                     // 测试 Redis 连接是否正常
// - Stats() *redis.PoolStats         // 获取连接池统计信息（用于监控）
// - InitRedisCluster() error         // 支持 Redis 集群模式
