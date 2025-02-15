package main

import "github.com/gin-gonic/gin"

func main() {
	err := initdatabase()
	if err != nil {
		panic(err)
	}

	Router := gin.Default()

	//登录前
	Router.POST("/user/register", Register) //注册
	Router.GET("/user/token", Login)        //登录（获取token）

	protectedRouter := Router.Group("/")
	protectedRouter.Use(IsLogin())
	//登陆后
	{
		protectedRouter.GET("/user/token/refresh", Refresh)              //刷新token（维持登录状态）
		protectedRouter.PUT("/user/password", Password)                  //修改用户密码
		protectedRouter.GET("/user/info/{user_id}", GetInfoById)         //获取用户信息
		protectedRouter.PUT("/user/info", Info)                          //修改用户信息
		protectedRouter.GET("/product/list", ShowList)                   //获取商品列表
		protectedRouter.POST("/book/search", SearchProduct)              //搜索商品
		protectedRouter.PUT("/product/addCart", AddCart)                 //加⼊购物⻋
		protectedRouter.GET("/product/cart", ShowCart)                   //获取购物⻋商品列表
		protectedRouter.GET("/product/info/{product_id}", ProductDetail) //获取商品详情
		protectedRouter.GET("/product/{type}", GetType)                  //获取相应标签的商品列表
		protectedRouter.GET("/comment/{product_id}", GetComment)         //获取商品的评论
		protectedRouter.POST("/comment/{product_id}", PutComment)        //给商品评论
		protectedRouter.DELETE("/comment/{comment_id}", DeleteComment)   //删除评论
		protectedRouter.PUT("/comment/{comment_id}", UpdateComment)      //更新评论
		protectedRouter.PUT("/comment/praise", Praise)                   //点赞点踩
		protectedRouter.POST("/operate/order", MakeOrder)                //下单
	}

	err = Router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
