package container

import (
	"ginApp/pkg/cache"
	"ginApp/pkg/database"
	"ginApp/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
)

var systemContainer *Container

// Container 服务容器
type Container struct {
	Db    *gorm.DB              //数据库
	Redis redis.UniversalClient //缓存
	Amqp  *amqp.Connection      //mq客户端
}

// GetContainer 获取容器实例
func GetContainer() *Container {
	return systemContainer
}

/*
Initialize
初始化基础容器
*/
func Initialize() {
	//数据库
	if err := database.InitPostgres(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	//缓存
	if err := cache.InitRedis(); err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	//队列
	if err := messaging.InitAmqp(); err != nil {
		log.Fatalf("初始化 消息队列 失败: %v", err)
	}

	//声明全局基类容器
	systemContainer = &Container{
		Db:    database.DB,
		Redis: cache.Client,
		Amqp:  messaging.Mq,
	}
}
