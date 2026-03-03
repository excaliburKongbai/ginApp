package router

import (
	"ginApp/internal/container"
	"ginApp/internal/handler/user"
)

// InitUserAccountRouter 初始化用户路由
func InitUserAccountRouter(r *RouterGroup, c *container.Container) {
	// 依赖注入：Repository -> Service -> Handler
	handler := user.New(c)

	// 公共接口（无需鉴权）
	r.Public.POST("/login", handler.Login)
	r.Public.POST("/register", handler.Register)
	r.Public.POST("/reset-password", handler.ResetPassword)

	// 需要鉴权的接口
	r.Protected.GET("/user/info", handler.Info)
	r.Protected.POST("/logout", handler.Logout)
}
