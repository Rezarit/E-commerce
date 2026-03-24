package bootstrap

import (
	"github.com/Rezarit/go-seckill-system/dao"
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/Rezarit/go-seckill-system/service"
	"github.com/Rezarit/go-seckill-system/service/consumers"
	"log"
)

func initDatabase() error {
	log.Println("[Bootstrap] 开始初始化数据库...")
	err := dao.InitDatabase()
	if err != nil {
		log.Printf("[Bootstrap] 数据库初始化失败: %v", err)
		return err
	}
	log.Println("[Bootstrap] 数据库初始化成功")
	return nil
}

func initConfig() error {
	log.Println("[Bootstrap] 开始加载配置...")
	config.InitConfig()
	log.Println("[Bootstrap] 配置加载成功")
	return nil
}

func initRedis() error {
	log.Println("[Bootstrap] 开始初始化Redis...")
	err := redis.InitRedis(config.GetRedisConfig())
	if err != nil {
		log.Printf("[Bootstrap] Redis初始化失败: %v", err)
		return err
	}
	log.Println("[Bootstrap] Redis初始化成功")
	return nil
}

func initMQ() error {
	log.Println("[Bootstrap] 开始初始化MQ...")

	// 获取MQ配置
	mqCfg := config.GetMQConfig()
	if mqCfg == nil || mqCfg.URL == "" {
		log.Println("[Bootstrap] MQ配置无效，跳过初始化")
		return nil
	}

	// 初始化
	if err := rabbitmq.InitRabbitMQ(mqCfg.URL, mqCfg.Queues); err != nil {
		log.Printf("[Bootstrap] MQ初始化失败: %v", err)
		return err
	}

	log.Println("[Bootstrap] MQ初始化成功")

	consumers.InitOrderConsumer()
	log.Println("[Bootstrap] 订单消费者已启动，在后台等待处理任务...")
	consumers.InitCartConsumer()
	log.Println("[Bootstrap] 购物车消费者已启动，在后台等待处理任务...")
	return nil
}

func initAllProductStock() error {
	log.Println("[Bootstrap] 开始初始化商品库存...")
	err := service.InitAllProductStock()
	if err != nil {
		log.Printf("[Bootstrap] 商品库存初始化失败: %v", err)
		return err
	}
	log.Println("[Bootstrap] 商品库存初始化成功")
	return nil
}

func Init() error {
	log.Println("开始应用初始化...")

	// 基础设施初始化
	if err := initConfig(); err != nil {
		return err
	}
	if err := initDatabase(); err != nil {
		return err
	}
	if err := initRedis(); err != nil {
		return err
	}
	if err := initMQ(); err != nil {
		return err
	}

	// 业务服务初始化
	if err := service.LoadLuaScripts(); err != nil {
		return err
	} // 脚本业务服务初始化
	service.InitService(redis.GetClient()) // 非脚本业务服务初始化

	// 缓存预热
	if err := initAllProductStock(); err != nil {
		return err
	}

	log.Println("应用初始化完成")
	return nil
}
