package test

import (
	"fmt"
)

type SettingKey[T any] string

// 定义所有配置键
var (
	AppName        = SettingKey[string]("app.name")
	RequestTimeout = SettingKey[int]("http.timeout.ms")
	DebugMode      = SettingKey[bool]("debug.mode")
)

// 配置存储结构
type Config struct {
	data map[string]interface{}
}

// 类型安全的 Get 方法
func Get[T any](c *Config, key SettingKey[T]) T {
	if val, ok := c.data[string(key)]; ok {
		if typedVal, ok := val.(T); ok {
			return typedVal
		}
	}
	var zero T
	return zero
}

// 类型安全的 Set 方法
func Set[T any](c *Config, key SettingKey[T], value T) {
	c.data[string(key)] = value
}

// 使用示例
func main() {
	config := &Config{
		data: make(map[string]interface{}),
	}

	// 设置配置
	Set(config, AppName, "MyApp")     // 只能设置 string
	Set(config, RequestTimeout, 3000) // 只能设置 int
	Set(config, DebugMode, true)      // 只能设置 bool

	// 获取配置（类型安全）
	name := Get(config, AppName)           // 自动推断为 string
	timeout := Get(config, RequestTimeout) // 自动推断为 int
	debug := Get(config, DebugMode)        // 自动推断为 bool

	fmt.Println(name, timeout, debug)
}
