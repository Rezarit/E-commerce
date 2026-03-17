package dao

import (
	"fmt"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

func InitDatabase() error {
	// 从配置中获取数据库配置
	dbCfg := config.GetDatabaseConfig()

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	log.Println("[数据库] 开始初始化数据库连接...")
	log.Printf("[数据库] 连接信息: %s@%s:%d/%s",
		dbCfg.User, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

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
	// 连接池配置（优化高并发场景）
	sqlDB.SetMaxOpenConns(500)              // 最大打开连接数：增加到500
	sqlDB.SetMaxIdleConns(200)              // 最大空闲连接数：增加到200
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