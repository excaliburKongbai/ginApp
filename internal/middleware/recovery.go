package middleware

import (
	"fmt"
	"ginApp/pkg/env"
	"ginApp/pkg/response"
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录日志 + 堆栈
				// 返回统一错误体
				errMsg := fmt.Sprintf("%s", err)
				if env.Active().IsPro() {
					//响应错误
					response.Error(c, 500, "服务器内部错误")
				} else {
					response.Error(c, 500, fmt.Sprintf(" 请求发生错误:%v", err))
				}
				//打印堆栈信息
				log.Println(
					fmt.Sprintf("\n\n -----------[error: %s]-----------\n\r  %s \n\r -----------[end]-----------", errMsg, string(debug.Stack())),
				)
			}
		}()
		c.Next()
	}
}
