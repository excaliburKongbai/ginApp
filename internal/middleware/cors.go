package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 什么是跨域？
// 浏览器的安全策略规定：前端页面（如 http://localhost:3000）
// 不能直接访问不同域名的后端接口（如 http://localhost:8080）
// CORS 中间件就是用来解决这个问题的
func CORS() gin.HandlerFunc {
	// gin.HandlerFunc 是一个函数类型，用于处理 HTTP 请求
	return func(c *gin.Context) {
		// 设置响应头，告诉浏览器允许跨域访问
		
		// Access-Control-Allow-Origin: * 表示允许任何域名访问
		// 生产环境建议改为具体的前端域名，如 "http://yourdomain.com"
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// 允许的请求头
		// Content-Type: 请求体类型（如 application/json）
		// Authorization: 用于传递 JWT token 等认证信息
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS 是浏览器的预检请求
		// 浏览器在发送真正的请求前，会先发送 OPTIONS 请求询问服务器是否允许跨域
		if c.Request.Method == "OPTIONS" {
			// AbortWithStatus 终止请求处理，直接返回 204 状态码
			// 204 表示成功但没有内容返回
			c.AbortWithStatus(204)
			return
		}

		// c.Next() 继续执行下一个中间件或处理器
		c.Next()
	}
}
