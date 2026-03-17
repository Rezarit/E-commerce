package route

import (
	"github.com/Rezarit/go-seckill-system/api"
	"github.com/Rezarit/go-seckill-system/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	Router := gin.Default()

	//登录前
	Router.POST("/user/register", api.Register)               //注册
	Router.POST("/user/login", api.Login)                     //用户登录
	Router.GET("/user/token/refresh", api.RefreshAccessToken) //刷新登录态

	// 商品相关
	Router.GET("/product/list", api.ShowProductList)           //获取商品列表
	Router.POST("/product/search", api.SearchProduct)          //搜索商品
	Router.GET("/product/info/:product_id", api.ProductDetail) //获取商品详情

	protectedRouter := Router.Group("/")
	protectedRouter.Use(middleware.LoginRequired())
	//登陆后
	{
		// 用户相关
		protectedRouter.PUT("/user/password", api.UpdateUserPassword)         //更新用户密码
		protectedRouter.GET("/user/info", api.GetUserInfoById)                //获取用户信息
		protectedRouter.PUT("/user/info", api.UpdateUserInfoByID)             //更新用户信息
		protectedRouter.POST("/user/register_merchant", api.RegisterMerchant) //注册商户

		// 商户相关
		roleMerchantRouter := protectedRouter.Group("/merchant")
		roleMerchantRouter.Use(middleware.MerchantRequired())
		roleMerchantRouter.POST("/product", api.CreatProduct)    //发布商品
		roleMerchantRouter.PUT("/product", api.UpdateProduct)    //更新商品
		roleMerchantRouter.DELETE("/product", api.DeleteProduct) //删除商品

		// 购物车相关
		protectedRouter.POST("/cart/add/:product_id", api.AddToCart)       //加⼊购物⻋
		protectedRouter.GET("/cart/list", api.ShowCart)                    //获取购物⻋商品列表
		protectedRouter.DELETE("/cart/remove/:product_id", api.RemoveCart) //从购物车移除商品

		// 订单相关
		protectedRouter.POST("/order/create", api.MakeOrder)             //下单
		protectedRouter.GET("/order/list", api.GetOrderList)             //获取订单列表
		protectedRouter.GET("/order/info/:order_id", api.GetOrderDetail) //获取订单详情
	}
	return Router
}
