package user

// Account 用户账号
import (
	"ginApp/internal/Dto/Response/user"
	"time"
)

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

func (User) TableComment() string {
	return "用户表"
}

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null;comment:'用户名称'" json:"username"`
	Password  string    `gorm:"size:255;not null;comment:'密码'" json:"-"` // json:"-" 不返回密码
	Nickname  string    `gorm:"size:50;comment:'昵称'" json:"nickname"`
	Email     string    `gorm:"uniqueIndex;size:100;comment:'邮箱'" json:"email"`
	Mobile    string    `gorm:"size:20;comment:'手机号码'" json:"mobile"`
	Avatar    string    `gorm:"size:255;comment:'头像'" json:"avatar"`
	Status    int8      `gorm:"default:1;comment:'状态:0=禁用,1=启用'" json:"status"` // 1:启用 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *user.UserResponse {
	return &user.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Mobile:    u.Mobile,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
