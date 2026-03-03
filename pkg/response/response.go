package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
	ErrCode *int        `json:"errcode,omitempty"`
}

const (
	ERROR   = -1
	SUCCESS = 0
)

// Result 成功的响应
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    data,
		Msg:     msg,
		ErrCode: nil,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "succeed", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "succeed", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// ErrResult 错误的响应
func ErrResult(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code,
		Data:    data,
		Msg:     msg,
		ErrCode: nil,
	})
}

// ErrWithDetailed 详细的错误响应
func ErrWithDetailed(data interface{}, message string, c *gin.Context) {
	ErrResult(ERROR, data, message, c)
}

// ErrWithMessage 错误的信息响应
func ErrWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

// ErrWithError 异常错误响应
func ErrWithError(err error, c *gin.Context) {
	ErrResult(ERROR, map[string]interface{}{}, fmt.Sprintf("异常错误:%s", err.Error()), c)
}

// ErrWithMessageByError 带错误描述的异常响应
func ErrWithMessageByError(message string, err error, c *gin.Context) {
	ErrResult(ERROR, map[string]interface{}{}, fmt.Sprintf("%s:%s", message, err.Error()), c)
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	Result(SUCCESS, data, "success", c)
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	Result(SUCCESS, data, msg, c)
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	Result(code, map[string]interface{}{}, msg, c)
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
