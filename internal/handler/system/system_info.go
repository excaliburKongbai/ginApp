package system

import (
	"ginApp/pkg/env"
	"ginApp/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Info 获取用户信息
func (h *Handler) Info(c *gin.Context) {
	// 从 token 或 session 中获取用户 ID（这里简化处理）
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.ErrWithMessage("缺少用户ID", c)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrWithMessage("用户ID格式错误", c)
		return
	}

	user, err := h.userService.GetUserInfo(c.Request.Context(), uint(userID))
	if err != nil {
		response.ErrWithMessageByError("获取用信息失败", err, c)
		return
	}

	response.OkWithDetailed(user, "success", c)
}

// Env 获取env版本
func (h *Handler) Env(c *gin.Context) {
	response.OkWithDetailed(gin.H{
		"env": env.Active().Value(),
	}, "success", c)
}

// Sms 发生短信
func (h *Handler) Sms(c *gin.Context) {}

// Ems 发送电子邮件短信
func (h *Handler) Ems(c *gin.Context) {}

// Notify 站内通知
func (h *Handler) Notify(c *gin.Context) {}
