package config

import (
	"github.com/spf13/viper"
	"log"
)

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret                 string `mapstructure:"secret"`
	AccessTokenExpireHour  int    `mapstructure:"access_token_expire_hour"`
	RefreshTokenExpireHour int    `mapstructure:"refresh_token_expire_hour"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

var Cfg GlobalConfig

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	// 解析到全局配置
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
}
