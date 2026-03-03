package messaging

import (
	"fmt"
	"ginApp/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Mq *amqp.Connection

/*
InitAmqp
初始化mq链接
*/
func InitAmqp() error {
	cfg := config.Cfg.Amqp
	var err error
	Mq, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	))
	if err != nil {
		return fmt.Errorf("连接rabbitmq失败: %w", err)
	}
	return nil
}

/*
CloseAmqp
关闭mq连接
*/
func CloseAmqp() error {
	return Mq.Close()
}
