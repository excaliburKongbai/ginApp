package core

import (
	"flag"
	"fmt"
	"ginApp/internal/container"
	"ginApp/internal/model/user"
	"log"
)

func DbInitialize() {
	// 解析命令行参数
	migrate := flag.Bool("migrate", false, "是否执行数据库迁移")
	flag.Parse()

	if *migrate {
		if err := authMigration(
			&user.User{}, //用户表
		); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
	}
}

// authMigration 迁移项目
func authMigration(models ...interface{}) error {
	DB := container.GetContainer().Db
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	log.Println("开始自动迁移数据库表结构...")
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println(fmt.Sprintf("数据库[%d]表结构迁移完成", len(models)))
	return nil
}
