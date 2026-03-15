package api

import (
	"github.com/Rezarit/go-seckill-system/api/common"
	"github.com/Rezarit/go-seckill-system/domain"
	"github.com/Rezarit/go-seckill-system/pkg/response"
	"github.com/Rezarit/go-seckill-system/service"
	"github.com/gin-gonic/gin"
	"log"
)

func Register(client *gin.Context) {
	var userRegisterRequest domain.UserRegisterRequest
	if isPass := common.BindRequest(client, &userRegisterRequest); !isPass {
		return
	}

	// 执行业务逻辑
	user, err := service.Register(userRegisterRequest)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "用户创建成功", domain.UserRegisterResponse{
		UserID: user.UserID,
	})
}

func Login(client *gin.Context) {
	var userLoginRequest domain.UserLoginRequest
	if isPass := common.BindRequest(client, &userLoginRequest); !isPass {
		return
	}

	// 执行业务逻辑
	userID, accessToken, refreshToken, err := service.Login(userLoginRequest)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "登录成功", domain.UserLoginResponse{
		UserID:       userID,
		Username:     userLoginRequest.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshAccessToken 刷新AT
func RefreshAccessToken(client *gin.Context) {
	// 读取Authorization请求头
	authHeader := client.GetHeader("Authorization")
	if authHeader == "" {
		log.Printf("[登录态验证失败] 请求头Authorization为空 | 请求路径：%s", client.FullPath())
		response.Fail(client, domain.ErrCodeTokenEmpty, "请先登录")
		return
	}

	refreshedAccessToken, err := service.RefreshAccessToken(authHeader)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "AT刷新成功", domain.UserRefreshATResponse{
		Token: refreshedAccessToken,
	})
}

func UpdateUserPassword(client *gin.Context) {
	//绑定数据
	var req domain.UserUpdatePasswordRequest
	isPass := common.BindRequest(client, &req)
	if !isPass {
		return
	}
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	//执行更新指令
	err := service.UpdateUserPassword(userID, req.Username, req.Password, req.NewPassword)
	if !common.HandleBusinessError(client, err) {
		return
	}

	//成功响应
	response.Success(client, "密码更新成功", nil)
}

func ParseUserID(client *gin.Context) int64 {
	userID, exists := client.Get("user_id")
	if !exists {
		response.Fail(client, domain.ErrCodeTokenInvalid, "用户未登录")
		return 0
	}
	return userID.(int64)
}

func UpdateUserInfoByID(client *gin.Context) {
	//获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}
	//绑定数据
	var user domain.User
	if isPass := common.BindRequest(client, &user); !isPass {
		return
	}

	// 执行业务逻辑
	err := service.UpdateUserInfoByID(userID, user)
	if !common.HandleBusinessError(client, err) {
		return
	}

	response.Success(client, "用户信息更新成功", nil)
}

func GetUserInfoById(client *gin.Context) {
	//获取用户ID
	userID := ParseUserID(client)
	if userID == 0 {
		return
	}

	//执行查询指令
	userInfo, err := service.GetUserInfoById(userID)
	if !common.HandleBusinessError(client, err) {
		return
	}

	//成功响应
	response.Success(client, "用户信息查询成功", userInfo)
}
