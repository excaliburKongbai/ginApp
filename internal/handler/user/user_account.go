package user

import (
	userRequest "ginApp/internal/Dto/Request/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req userRequest.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req userRequest.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "注册成功",
		"data": user,
	})
}

// Info 获取用户信息
func (h *Handler) Info(c *gin.Context) {
	// 从 token 或 session 中获取用户 ID（这里简化处理）
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少用户ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式错误"})
		return
	}

	user, err := h.userService.GetUserInfo(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": user,
	})
}

// Logout 用户登出
func (h *Handler) Logout(c *gin.Context) {
	// TODO: 清除 token 或 session
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登出成功",
	})
}

// ResetPassword 重置密码
func (h *Handler) ResetPassword(c *gin.Context) {
	var req userRequest.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码重置成功",
	})
}
