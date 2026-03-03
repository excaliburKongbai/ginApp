package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
	ErrCode *int        `json:"errcode,omitempty"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

// Result 成功的响应
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
		nil,
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
		code,
		data,
		msg,
		nil,
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
