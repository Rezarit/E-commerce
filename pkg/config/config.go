package config

import (
	"github.com/Rezarit/go-seckill-system/pkg/logger"
	"github.com/spf13/viper"
)

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret                 string `mapstructure:"secret"`
	AccessTokenExpireHour  int    `mapstructure:"access_token_expire_hour"`
	RefreshTokenExpireHour int    `mapstructure:"refresh_token_expire_hour"`
}

// MQConfig RabbitMQ配置结构体
type MQConfig struct {
	URL    string            `mapstructure:"url"`
	Queues map[string]string `mapstructure:"queues"`
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Addr     string `yaml:"addr"` // Redis地址: "localhost:6379"
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Database DatabaseConfig `mapstructure:"database"`
	RabbitMQ MQConfig       `mapstructure:"rabbitmq"`
	Logger   logger.Config  `mapstructure:"logger"`
}

var Cfg GlobalConfig

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 解析到全局配置
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}

func GetRedisConfig() *RedisConfig {
	redisCfg := &RedisConfig{}
	err := viper.UnmarshalKey("redis", redisCfg)
	if err != nil {
		return nil
	}
	return redisCfg
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() *DatabaseConfig {
	dbCfg := &DatabaseConfig{}

	if !viper.IsSet("database") {
		panic("配置文件中未找到database配置，请在config.yaml中配置数据库信息")
	}

	err := viper.UnmarshalKey("database", dbCfg)
	if err != nil {
		panic(err)
	}

	// 验证必填字段
	if dbCfg.Host == "" || dbCfg.User == "" || dbCfg.DBName == "" {
		panic("数据库配置不完整，请检查config.yaml文件")
	}

	return dbCfg
}

func GetMQConfig() *MQConfig {
	// 实现配置获取逻辑
	mqCfg := &MQConfig{}
	err := viper.UnmarshalKey("rabbitmq", mqCfg)
	if err != nil {
		return nil
	}
	return mqCfg
}
