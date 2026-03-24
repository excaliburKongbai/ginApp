package system

import (
	"ginApp/internal/container"
	"ginApp/internal/repository/user"
	userService "ginApp/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ InterfaceSystem = (*Handler)(nil)

type InterfaceSystem interface {
	//基础信息
	Info(c *gin.Context)
	//环境变量
	Env(c *gin.Context)
	//短信
	Sms(c *gin.Context)
	//电子邮件
	Ems(c *gin.Context)
	//通知
	Notify(c *gin.Context)
}

type Handler struct {
	userService userService.Service
	Db          *gorm.DB
	Cache       redis.UniversalClient
}

func New(c *container.Container) *Handler {
	return &Handler{
		userService: userService.NewService(user.NewRepository(c.Db)),
		Db:          c.Db,
		Cache:       c.Redis,
	}
}
