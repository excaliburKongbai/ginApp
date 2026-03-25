package router

import (
	_ "ginApp/docs" // 这里改成你的 module 路径 + /docs
	"ginApp/internal/container"
	"ginApp/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	//文档路由
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	return Router
}
