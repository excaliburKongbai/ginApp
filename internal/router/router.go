package router

import (
	"fmt"
	"ginApp/internal/container"
	"ginApp/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RouterGroup 路由分组
type RouterGroup struct {
	Public    *gin.RouterGroup //非鉴权
	Protected *gin.RouterGroup //鉴权
}

// 初始化总路由
func Routers(container *container.Container) *gin.Engine {
	Router := gin.New()

	//使用中间件
	Router.Use(middleware.Recovery()) //错误收集
	Router.Use(middleware.CORS())     //响应头
	Router.Use(middleware.Logger())   //响应时间

	apiRouter := Router.Group("/api")
	//公共接口
	routerGroup := &RouterGroup{
		apiRouter.Group("/"),
		apiRouter.Group("/auth"),
	}

	//用户接口
	InitUserAccountRouter(routerGroup, container)
	//系统入口
	InitSystemRouter(routerGroup, container)

	defer func() {
		fmt.Println("执行完毕了接口")
	}()

	return Router
}
