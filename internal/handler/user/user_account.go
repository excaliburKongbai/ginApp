package user

import (
	userRequest "ginApp/internal/Dto/Request/user"
	"ginApp/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req userRequest.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		response.ErrWithError(err, c)
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.ErrWithMessageByError("用户登陆失败", err, c)
		return
	}

	response.OkWithDetailed(gin.H{
		"user":  user,
		"token": token,
	}, "登陆成功", c)
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req userRequest.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrWithError(err, c)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		response.ErrWithMessageByError("用户注册失败", err, c)
		return
	}
	response.OkWithDetailed(user, "注册成功", c)
}

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

// Logout 用户登出
func (h *Handler) Logout(c *gin.Context) {
	// TODO: 清除 token 或 session
	response.OkWithMessage("登出成功", c)
}

// ResetPassword 重置密码
func (h *Handler) ResetPassword(c *gin.Context) {
	var req userRequest.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrWithError(err, c)
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), &req); err != nil {
		response.ErrWithMessageByError("修改密码失败", err, c)
		return
	}

	response.OkWithMessage("密码重置成功", c)
}
