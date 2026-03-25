package router

import (
	"ginApp/internal/container"
	"ginApp/internal/handler/system"
)

// InitSystemRouter 初始化系统路由
func InitSystemRouter(r *RouterGroup, c *container.Container) {
	// 依赖注入：Repository -> Service -> Handler
	handler := system.New(c)

	//系统信息
	r.Public.GET("/info", handler.Info)

	//env环境
	r.Public.GET("/env", handler.Env)
}
