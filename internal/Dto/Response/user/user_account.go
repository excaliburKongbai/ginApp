package user

import "time"

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Avatar    string    `json:"avatar"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Avatar    string    `json:"avatar"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
