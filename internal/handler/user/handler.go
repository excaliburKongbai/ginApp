package user

import (
	"ginApp/internal/container"
	"ginApp/internal/repository/user"
	userService "ginApp/internal/service/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	userService userService.Service
	Db          *gorm.DB
	Cache       redis.UniversalClient
}

func New(c *container.Container) *Handler {
	return &Handler{
		userService.NewService(user.NewRepository(c.Db)),
		c.Db,
		c.Redis,
	}
}
