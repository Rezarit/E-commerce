package dao

import (
	"github.com/Rezarit/E-commerce/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

func InitDatabase() error {
	dsn := "root:fzfz1314@tcp(127.0.0.1:3306)/e_commerce?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化 GORM 连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 开发环境打印 SQL 日志，方便调试；生产环境可改为 logger.Silent
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用默认事务
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	// 获取底层 sql.DB，配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	// 连接池配置
	sqlDB.SetMaxOpenConns(20)                  // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接存活时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 连接空闲超时

	// 自动迁移表
	err = DB.AutoMigrate(&domain.User{})
	if err != nil {
		return err
	}

	return nil
}
