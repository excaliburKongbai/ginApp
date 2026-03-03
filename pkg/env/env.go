// Package env 提供环境管理功能，用于区分不同的运行环境（开发、测试、预发布、生产）
package env

import (
	"flag"
	"fmt"
	"strings"
)

var (
	// active 当前激活的环境，通过命令行参数 -env 设置
	active Environment
	
	// dev 开发环境实例
	dev Environment = &environment{value: "dev"}
	
	// fat 测试环境实例 (Feature Acceptance Test)
	fat Environment = &environment{value: "fat"}
	
	// uat 预上线环境实例 (User Acceptance Test)
	uat Environment = &environment{value: "uat"}
	
	// pro 生产环境实例 (Production)
	pro Environment = &environment{value: "pro"}
)

// 编译时检查：确保 environment 实现了 Environment 接口
var _ Environment = (*environment)(nil)

// Environment 环境接口，定义了环境相关的方法
type Environment interface {
	// Value 返回环境的字符串值（"dev", "fat", "uat", "pro"）
	Value() string
	
	// IsDev 判断是否为开发环境
	IsDev() bool
	
	// IsFat 判断是否为测试环境
	IsFat() bool
	
	// IsUat 判断是否为预上线环境
	IsUat() bool
	
	// IsPro 判断是否为生产环境
	IsPro() bool
	
	// t 私有方法，防止外部实现此接口
	t()
}

// environment 环境的具体实现
type environment struct {
	value string // 环境名称
}

// Value 返回当前环境的字符串值
func (e *environment) Value() string {
	return e.value
}

// IsDev 判断当前环境是否为开发环境
func (e *environment) IsDev() bool {
	return e.value == "dev"
}

// IsFat 判断当前环境是否为测试环境
func (e *environment) IsFat() bool {
	return e.value == "fat"
}

// IsUat 判断当前环境是否为预上线环境
func (e *environment) IsUat() bool {
	return e.value == "uat"
}

// IsPro 判断当前环境是否为生产环境
func (e *environment) IsPro() bool {
	return e.value == "pro"
}

// t 私有方法，确保只有包内部可以创建 Environment 实例
func (e *environment) t() {}

// init 包初始化函数，在程序启动时自动执行
// 从命令行参数读取环境配置，并设置当前激活的环境
func init() {
	// 定义命令行参数 -env，用于指定运行环境
	env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n")
	
	// 解析命令行参数
	flag.Parse()

	// 根据输入的环境参数设置当前激活的环境
	// 转换为小写并去除空格，确保参数格式统一
	switch strings.ToLower(strings.TrimSpace(*env)) {
	case "dev":
		active = dev
	case "fat":
		active = fat
	case "uat":
		active = uat
	case "pro":
		active = pro
	default:
		// 如果没有指定环境或环境参数非法，默认使用 fat（测试环境）
		active = fat
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'fat' will be used.")
	}
}

// Active 当前配置的env
func Active() Environment {
	return active
}
