package service

import (
	"os"

	"github.com/Rezarit/go-seckill-system/pkg/logger"
	myredis "github.com/Rezarit/go-seckill-system/pkg/redis"
	"github.com/go-redis/redis/v8"
)

const (
	ScriptAddToCart   = "scripts/lua/add_to_cart.lua"
	ScriptDeductStock = "scripts/lua/deduct_stock.lua"
	ScriptAddBackStock = "scripts/lua/add_back_stock.lua"
)

type CartService struct {
	client    *redis.Client
	luaScript *redis.Script
}

type StockDeductService struct {
	client        *redis.Client
	luaScript     *redis.Script // 扣减脚本
	addBackScript *redis.Script // 加回脚本（补偿）
}

type CacheService struct {
	client *redis.Client
}

type OrderService struct {
	client *redis.Client
}

var (
	cartService        *CartService
	stockDeductService *StockDeductService
	cacheService       *CacheService
	Order              *OrderService
)

type MessageHandler func(body []byte) error

// LoadLuaScripts 初始化lua脚本
func LoadLuaScripts() error {
	logger.Sugar.Info("[Service] 开始初始化服务实例...")

	// 初始化添加到购物车服务
	var err error
	cartService, err = NewLuaScriptService(
		ScriptAddToCart,
		func(client *redis.Client, script *redis.Script) *CartService {
			return &CartService{client: client, luaScript: script}
		})
	if err != nil {
		logger.Sugar.Errorf("[Service] 购物车服务初始化失败: %v", err)
		return err
	}
	logger.Sugar.Info("[Service] 购物车服务初始化完成")

	// 初始化库存减扣服务（需要两个脚本：扣减 + 加回补偿）
	deductScript, err := loadLuaScript(ScriptDeductStock)
	if err != nil {
		logger.Sugar.Errorf("[Service] 扣减脚本加载失败: %v", err)
		return err
	}
	addBackScript, err := loadLuaScript(ScriptAddBackStock)
	if err != nil {
		logger.Sugar.Errorf("[Service] 加回脚本加载失败: %v", err)
		return err
	}
	stockDeductService = &StockDeductService{
		client:        myredis.GetClient(),
		luaScript:     redis.NewScript(deductScript),
		addBackScript: redis.NewScript(addBackScript),
	}
	logger.Sugar.Info("[Service] 库存减扣服务初始化完成")

	logger.Sugar.Info("[Service] 所有服务实例初始化完成")
	return nil
}

// NewLuaScriptService 通用Lua脚本服务工厂函数
func NewLuaScriptService[T any](scriptPath string, constructor func(client *redis.Client, script *redis.Script) T) (T, error) {
	var zero T

	scriptContent, err := loadLuaScript(scriptPath)
	if err != nil {
		return zero, err
	}

	service := constructor(myredis.GetClient(), redis.NewScript(scriptContent))
	return service, nil
}

// loadLuaScript 从文件加载Lua脚本
func loadLuaScript(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Sugar.Errorf("加载Lua脚本失败: %v", err)
		return "", err
	}
	return string(content), nil
}

// InitService 初始化非脚本服务
func InitService(redisClient *redis.Client) {
	logger.Sugar.Info("[Service] 开始初始化服务实例...")

	cacheService = &CacheService{client: redisClient}
	Order = &OrderService{client: redisClient}

	logger.Sugar.Info("[Service] 服务实例初始化成功")
}
