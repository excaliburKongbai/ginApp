package response

import "github.com/gin-gonic/gin"

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(400, Response{
		Code: 400,
		Msg:  msg,
	})
}

// Unauthorized 401 未授权
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(401, Response{
		Code: 401,
		Msg:  msg,
	})
}

// NotFound 404 未找到
func NotFound(c *gin.Context, msg string) {
	c.JSON(404, Response{
		Code: 404,
		Msg:  msg,
	})
}
