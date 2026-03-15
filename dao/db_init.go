package dao

import (
	"github.com/Rezarit/go-seckill-system/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

func InitDatabase() error {
	log.Println("[数据库] 开始初始化数据库连接...")
	dsn := "root:fzfz1314@tcp(127.0.0.1:3306)/e_commerce?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化 GORM 连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Printf("[数据库] 初始化数据库连接失败: %v", err)
		return err
	}
	log.Println("[数据库] 数据库连接初始化成功")

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
	log.Println("[数据库] 开始自动迁移数据库表...")
	err = DB.AutoMigrate(
		&domain.User{},
		&domain.Merchant{},
		&domain.MerchantApplication{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
	)
	if err != nil {
		log.Printf("[数据库] 自动迁移数据库表失败: %v", err)
		return err
	}
	log.Println("[数据库] 数据库表自动迁移成功")

	return nil
}
