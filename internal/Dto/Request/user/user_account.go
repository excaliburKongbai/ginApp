package user

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"email"`
	Mobile   string `json:"mobile"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	Code        string `json:"code" binding:"required"` // 验证码
}
