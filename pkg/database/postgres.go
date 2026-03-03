package database

import (
	"fmt"
	"log"
	"time"

	"ginApp/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
//
// 为什么使用全局变量？
//  1. 数据库连接池应该在整个应用生命周期内只创建一次
//  2. 连接池是线程安全的，可以在多个 goroutine 中并发使用
//  3. 避免重复创建连接，提高性能
//
// 注意：
//   - DB 不是单个连接，而是一个连接池
//   - GORM 会自动管理连接的获取和释放
//   - 不需要手动关闭连接，只需在应用退出时关闭连接池
var DB *gorm.DB

// InitPostgres 初始化 PostgreSQL 数据库连接池
//
// 这个函数在应用启动时调用一次（在 main.go 中）
// 它会：
//  1. 读取配置文件中的数据库配置
//  2. 构建数据库连接字符串（DSN）
//  3. 创建数据库连接池
//  4. 配置连接池参数（最大连接数、空闲连接数等）
//  5. 测试连接是否正常
//
// 返回:
//   - error: 如果连接失败返回错误，成功返回 nil
func InitPostgres() error {
	// ==================== 第一步：读取配置 ====================
	// 从配置文件中读取数据库配置
	// config.Cfg 是全局配置对象，在 main.go 中通过 config.Load() 加载
	cfg := config.Cfg.Database

	// ==================== 第二步：构建 DSN（数据源名称）====================
	// DSN (Data Source Name) 是数据库连接字符串
	// 格式: host=主机 port=端口 user=用户名 password=密码 dbname=数据库名 sslmode=SSL模式
	//
	// 参数说明：
	//   - host: 数据库服务器地址（localhost 或 IP 地址）
	//   - port: 数据库端口（PostgreSQL 默认 5432）
	//   - user: 数据库用户名
	//   - password: 数据库密码
	//   - dbname: 要连接的数据库名称
	//   - sslmode: SSL 连接模式
	//       * disable: 不使用 SSL（开发环境）
	//       * require: 必须使用 SSL（生产环境）
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// ==================== 第三步：创建数据库连接 ====================
	var err error
	// gorm.Open 创建数据库连接池
	// 参数：
	//   1. postgres.Open(dsn): 指定使用 PostgreSQL 驱动和连接字符串
	//   2. &gorm.Config: GORM 配置选项
	//
	// Logger 配置：
	//   - logger.Silent: 不打印任何 SQL（生产环境推荐）
	//   - logger.Error: 只打印错误 SQL
	//   - logger.Warn: 打印慢查询和错误
	//   - logger.Info: 打印所有 SQL（开发环境推荐，方便调试）
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印所有 SQL 语句
	})
	if err != nil {
		// 连接失败，返回错误
		// %w 是 Go 1.13+ 的错误包装语法，可以保留原始错误信息
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// ==================== 第四步：配置连接池参数 ====================
	// 获取底层的 *sql.DB 对象
	// GORM 是基于 database/sql 包的 ORM，底层使用标准库的连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中的最大连接数
	// 空闲连接是已经建立但暂时没有使用的连接
	// 保留空闲连接可以提高性能，避免频繁创建和销毁连接
	// 建议值：5-10（根据实际负载调整）
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置数据库的最大打开连接数
	// 包括正在使用的连接和空闲连接
	// 限制最大连接数可以：
	//   1. 防止数据库压力过大
	//   2. 避免耗尽数据库连接资源
	//   3. 控制应用的资源使用
	// 建议值：10-100（根据数据库配置和应用负载调整）
	// 注意：不要设置太大，PostgreSQL 默认最大连接数是 100
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置连接的最大生命周期
	// 超过这个时间的连接会被关闭并重新创建
	// 好处：
	//   1. 避免长时间连接导致的问题（如网络中断、数据库重启）
	//   2. 定期刷新连接，保持连接健康
	//建议值：5-30 分钟
	sqlDB.SetConnMaxLifetime(time.Hour)

	// ==================== 第五步：测试连接 ====================
	// Ping() 测试数据库连接是否正常
	// 如果数据库无法访问，这里会返回错误
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Println("PostgreSQL 连接成功")
	return nil
}

// 未来可以添加的功能：
// - Close() error                    // 关闭数据库连接池（应用退出时调用）
// - Ping() error                     // 测试数据库连接是否正常
// - Stats() sql.DBStats              // 获取连接池统计信息（用于监控）
// - InitMySQL() error                // 支持 MySQL 数据库
// - InitSQLite() error               // 支持 SQLite 数据库（用于测试）
