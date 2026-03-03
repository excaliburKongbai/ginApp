package core

import (
	"ginApp/internal/config"
	"ginApp/internal/container"
	"ginApp/internal/router"
	"time"
)

// RunServer 运行服务
func RunServer() {
	//初始化路由参数
	Router := router.Routers(container.GetContainer())

	//启动服务
	initServer(config.Cfg.Server.Port, Router, 10*time.Minute, 10*time.Minute)
}
