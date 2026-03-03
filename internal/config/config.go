package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用的总配置结构
// 包含服务器、数据库、Redis、JWT 等所有配置项
type Config struct {
	Server   ServerConfig   `yaml:"server"` // `yaml:"server"` 表示对应 YAML 文件中的 server 字段
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Amqp     AmqpConfig     `yaml:"amqp"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"` // 服务器监听端口，如 8080
	Mode string `yaml:"mode"` // 运行模式: debug(开发), release(生产), test(测试)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`           // 数据库主机地址
	Port         string `yaml:"port"`           // 数据库端口
	User         string `yaml:"user"`           // 数据库用户名
	Password     string `yaml:"password"`       // 数据库密码
	DBName       string `yaml:"dbname"`         // 数据库名称
	SSLMode      string `yaml:"sslmode"`        // SSL 模式: disable(关闭), require(需要)
	MaxIdleConns int    `yaml:"max_idle_conns"` // 最大空闲连接数
	MaxOpenConns int    `yaml:"max_open_conns"` // 最大打开连接数
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`      // Redis 主机地址
	Port     string `yaml:"port"`      // Redis 端口
	Password string `yaml:"password"`  // Redis 密码
	DB       int    `yaml:"db"`        // Redis 数据库编号 (0-15)
	PoolSize int    `yaml:"pool_size"` // 连接池大小
}

// JWTConfig JWT (JSON Web Token) 配置
// JWT 用于用户身份认证，生成登录令牌
type JWTConfig struct {
	Secret      string `yaml:"secret"`       // JWT 签名密钥，用于加密和验证 token
	ExpireHours int    `yaml:"expire_hours"` // token 过期时间（小时）
}

// AmqpConfig 异步队列
type AmqpConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Cfg 全局配置实例
// 应用启动时加载配置文件，之后所有模块都可以通过 config.Cfg 访问配置
var Cfg *Config

// Load 从 YAML 文件加载配置
// path: 配置文件路径，如 "config/config.yaml"
func Load(path string) error {
	// os.ReadFile 读取文件内容到字节数组
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 创建一个空的 Config 结构体
	Cfg = &Config{}

	// yaml.Unmarshal 将 YAML 格式的字节数组解析到 Config 结构体
	// 这个过程叫做"反序列化"，把文本格式转换成 Go 的数据结构
	if err := yaml.Unmarshal(data, Cfg); err != nil {
		return err
	}

	log.Println("配置加载成功")
	return nil
}
