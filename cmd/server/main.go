package main

import (
	"fmt"
	"ginApp/core"
	"ginApp/internal/container"
	"log"

	"ginApp/internal/config"
	"github.com/gin-gonic/gin"
)

// main 应用程序入口
// 这是整个应用的启动点，负责初始化所有组件并启动服务器
// @title Gin App API
// @version 1.0
// @description API 文档
// @host

// @BasePath /api
func main() {
	defer systemOut()
	// ==================== 第一步：加载配置 ====================
	// 从 config/config.yaml 读取所有配置（数据库、Redis、服务器端口等）
	if err := config.Load("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// ==================== 第二步：设置 Gin 运行模式 ====================
	// debug: 开发模式，会打印详细日志和错误堆栈
	// release: 生产模式，性能更好，日志更少
	gin.SetMode(config.Cfg.Server.Mode)

	// ==================== 第三步：初始化基础设施层（最底层）====================
	// 基础设施层包括数据库、Redis、消息队列等中间件
	// 这些连接在应用启动时建立一次，整个应用生命周期内复用
	initializeSystem()

	// ==================== 第三步：启动服务====================
	core.RunServer()
}

/*
*
initializeSystem
初始化基类服务容器
*/
func initializeSystem() {
	//初始化服务基类
	container.Initialize()
	//db初始化
	core.DbInitialize()

}

// 系统结束
func systemOut() {
	if r := recover(); r != nil {
		log.Println(fmt.Sprintf("系统发生错误:%s", r))
	}
}
