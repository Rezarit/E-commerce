package bootstrap

import (
	"github.com/Rezarit/go-seckill-system/internal/dao"
	service2 "github.com/Rezarit/go-seckill-system/internal/service"
	consumers2 "github.com/Rezarit/go-seckill-system/internal/service/consumers"
	"github.com/Rezarit/go-seckill-system/pkg/config"
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/Rezarit/go-seckill-system/pkg/rabbitmq"
	"github.com/Rezarit/go-seckill-system/pkg/redis"
)

func initLogger() error {
	if err := logger.InitLogger(&config.Cfg.Logger); err != nil {
		return err
	}
	logger.Sugar.Info("[Bootstrap] 日志系统初始化成功")
	return nil
}

func initDatabase() error {
	logger.Sugar.Info("[Bootstrap] 开始初始化数据库...")
	err := dao.InitDatabase()
	if err != nil {
		logger.Sugar.Errorf("[Bootstrap] 数据库初始化失败: %v", err)
		return err
	}
	logger.Sugar.Info("[Bootstrap] 数据库初始化成功")
	return nil
}

func initConfig() error {
	config.InitConfig()
	return nil
}

func initRedis() error {
	logger.Sugar.Info("[Bootstrap] 开始初始化Redis...")
	err := redis.InitRedis(config.GetRedisConfig())
	if err != nil {
		logger.Sugar.Errorf("[Bootstrap] Redis初始化失败: %v", err)
		return err
	}
	logger.Sugar.Info("[Bootstrap] Redis初始化成功")
	return nil
}

func initMQ() error {
	logger.Sugar.Info("[Bootstrap] 开始初始化MQ...")

	// 获取MQ配置
	mqCfg := config.GetMQConfig()
	if mqCfg == nil || mqCfg.URL == "" {
		logger.Sugar.Info("[Bootstrap] MQ配置无效，跳过初始化")
		return nil
	}

	// 初始化
	if err := rabbitmq.InitRabbitMQ(mqCfg.URL, mqCfg.Queues); err != nil {
		logger.Sugar.Errorf("[Bootstrap] MQ初始化失败: %v", err)
		return err
	}

	logger.Sugar.Info("[Bootstrap] MQ初始化成功")
	return nil
}

func initConsumer() {
	consumers2.InitOrderConsumer()
	logger.Sugar.Info("[Bootstrap] 订单消费者已启动，在后台等待处理任务...")
	consumers2.InitCartConsumer()
	logger.Sugar.Info("[Bootstrap] 购物车消费者已启动，在后台等待处理任务...")
	consumers2.InitDLXConsumer()
	logger.Sugar.Info("[Bootstrap] 死信消费者已启动，在后台等待处理任务...")
}

func initAllProductStock() error {
	logger.Sugar.Info("[Bootstrap] 开始初始化商品库存...")
	err := service2.InitAllProductStock()
	if err != nil {
		logger.Sugar.Errorf("[Bootstrap] 商品库存初始化失败: %v", err)
		return err
	}
	logger.Sugar.Info("[Bootstrap] 商品库存初始化成功")
	return nil
}

func Init() error {
	// 先加载配置
	if err := initConfig(); err != nil {
		return err
	}

	// 初始化日志
	if err := initLogger(); err != nil {
		return err
	}

	logger.Sugar.Info("开始应用初始化...")

	// 基础设施初始化
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
	if err := service2.LoadLuaScripts(); err != nil {
		return err
	} // 脚本业务服务初始化
	service2.InitService(redis.GetClient()) // 非脚本业务服务初始化
	initConsumer()                          // 初始化消费者

	// 缓存预热
	if err := initAllProductStock(); err != nil {
		return err
	}

	logger.Sugar.Info("应用初始化完成")
	return nil
}
