package main

import (
	"github.com/Rezarit/E-commerce/api"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := dao.Initdatabase()
	if err != nil {
		panic(err)
	}

	Router := gin.Default()

	//登录前
	Router.POST("/user/register", api.Register) //注册
	Router.POST("/user/token", api.Login)       //登录（获取token）

	protectedRouter := Router.Group("/")
	protectedRouter.Use(service.IsLogin())
	//登陆后
	{
		//登录相关
		protectedRouter.GET("/user/token/refresh", api.Refresh) //刷新token（维持登录状态）

		//用户相关
		protectedRouter.PUT("/user/password", api.Password)          //修改用户密码
		protectedRouter.GET("/user/info/{user_id}", api.GetInfoById) //获取用户信息
		protectedRouter.PUT("/user/info", api.Info)                  //修改用户信息

		//商品相关
		protectedRouter.GET("/product/list", api.ShowList)                   //获取商品列表
		protectedRouter.POST("/book/search/{product_id}", api.SearchProduct) //搜索商品
		protectedRouter.GET("/product/info/{product_id}", api.ProductDetail) //获取商品详情
		protectedRouter.GET("/product/{type}", api.GetType)                  //获取相应标签的商品列表

		//购物车相关
		protectedRouter.PUT("/product/addCart/{product_id}", api.AddCart) //加⼊购物⻋
		protectedRouter.GET("/product/cart", api.ShowCart)                //获取购物⻋商品列表
		protectedRouter.POST("/operate/order", api.MakeOrder)             //下单

		//评论相关
		protectedRouter.GET("/comment/{product_id}", api.GetComment)       //获取商品的评论
		protectedRouter.POST("/comment/{product_id}", api.PutComment)      //给商品评论
		protectedRouter.DELETE("/comment/{comment_id}", api.DeleteComment) //删除评论
		protectedRouter.PUT("/comment/{comment_id}", api.UpdateComment)    //更新评论
		protectedRouter.PUT("/comment/praise", api.Praise)                 //点赞点踩
	}

	err = Router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
