package route

import (
	"github.com/Rezarit/E-commerce/api"
	"github.com/Rezarit/E-commerce/api/middleware"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	err := dao.InitDatabase()
	if err != nil {
		panic(err)
	}

	Router := gin.Default()

	//登录前
	Router.POST("/user/register", api.Register)               //注册
	Router.POST("/user/token", api.Login)                     //登录
	Router.GET("/user/token/refresh", api.RefreshAccessToken) //刷新登录态

	protectedRouter := Router.Group("/")
	protectedRouter.Use(middleware.LoginRequired())
	//登陆后
	{
		//用户相关
		protectedRouter.PUT("/user/update-password", api.UpdateUserPassword) //更新用户密码
		protectedRouter.GET("/user/info/{user_id}", api.GetUserInfoById)     //获取用户信息
		protectedRouter.PUT("/user/info", api.UpdateUserInfoByID)            //更新用户信息

		////商品相关
		//protectedRouter.GET("/product/list", api.ShowList)                   //获取商品列表
		//protectedRouter.POST("/book/search/{product_id}", api.SearchProduct) //搜索商品
		//protectedRouter.GET("/product/info/{product_id}", api.ProductDetail) //获取商品详情
		//protectedRouter.GET("/product/{type}", api.GetType)                  //获取相应标签的商品列表

		////购物车相关
		//protectedRouter.PUT("/product/addCart/{product_id}", api.AddCart) //加⼊购物⻋
		//protectedRouter.GET("/product/cart", api.ShowCart)                //获取购物⻋商品列表
		//protectedRouter.POST("/operate/order", api.MakeOrder)             //下单

		////评论相关
		//protectedRouter.GET("/comment/{product_id}", api.GetComment)       //获取商品的评论
		//protectedRouter.POST("/comment/{product_id}", api.PutComment)      //给商品评论
		//protectedRouter.DELETE("/comment/{comment_id}", api.DeleteComment) //删除评论
		//protectedRouter.PUT("/comment/{comment_id}", api.UpdateComment)    //更新评论
		//protectedRouter.PUT("/comment/praise", api.Praise)                 //点赞点踩
	}
	return Router
}
