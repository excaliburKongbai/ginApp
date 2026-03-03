package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
// 记录每个 HTTP 请求的详细信息，方便调试和监控
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// time.Now() 记录请求开始时间
		start := time.Now()
		
		// 获取请求的路径，如 /api/v1/users/123
		path := c.Request.URL.Path
		
		// 获取请求方法，如 GET, POST, PUT, DELETE
		method := c.Request.Method

		// c.Next() 执行实际的请求处理
		// 这之后的代码会在请求处理完成后执行
		c.Next()

		// time.Since(start) 计算请求处理耗时
		latency := time.Since(start)
		
		// c.Writer.Status() 获取响应状态码
		// 200 表示成功，404 表示未找到，500 表示服务器错误等
		status := c.Writer.Status()

		// 打印日志，格式: [方法] 路径 状态码 耗时
		// 例如: [GET] /api/v1/users/1 200 5.2ms
		log.Printf("[%s] %s %d %v", method, path, status, latency)
	}
}
