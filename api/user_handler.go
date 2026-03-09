package api

import (
	"errors"
	"github.com/Rezarit/E-commerce/api/auth"
	"github.com/Rezarit/E-commerce/api/common"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/response"
	"github.com/Rezarit/E-commerce/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Register(client *gin.Context) {
	var userRegisterRequest domain.UserRegisterRequest
	if isPass := common.BindRequest(client, &userRegisterRequest); !isPass {
		return
	}

	// 执行业务逻辑
	user, err := service.Register(userRegisterRequest)
	if err != nil {
		// 处理业务错误
		var bizErr *domain.BusinessError
		if !errors.As(err, &bizErr) {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}
		response.Fail(client, bizErr.Code, bizErr.Msg)
		return
	}

	response.Success(client, "用户创建成功", domain.UserRegisterResponse{
		Username: user.Username,
	})
}

func Login(client *gin.Context) {
	var userLoginRequest domain.UserLoginRequest
	if isPass := common.BindRequest(client, &userLoginRequest); !isPass {
		return
	}

	// 执行业务逻辑
	accessToken, refreshToken, err := service.Login(userLoginRequest)
	if err != nil {
		// 处理业务错误
		var bizErr *domain.BusinessError
		if !errors.As(err, &bizErr) {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}

		response.Fail(client, bizErr.Code, bizErr.Msg)
		return
	}

	response.Success(client, "登录成功", domain.UserLoginResponse{
		Username:     userLoginRequest.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshAccessToken 刷新AT
func RefreshAccessToken(client *gin.Context) {
	authHeader := auth.GetAuthHeader(client)
	refreshedAccessToken, err := service.RefreshAccessToken(authHeader)
	if err != nil {
		var bizErr *domain.BusinessError
		if !errors.As(err, &bizErr) {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}
		response.Fail(client, bizErr.Code, bizErr.Msg)
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

	//执行更新指令
	err := service.UpdateUserPassword(req.Username, req.Password, req.NewPassword)
	if err != nil {
		var bizErr *domain.BusinessError
		if !errors.As(err, &bizErr) {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}
		response.Fail(client, bizErr.Code, bizErr.Msg)
		return
	}

	//成功响应
	response.Success(client, "密码更新成功", nil)
}

func ParseUserID(client *gin.Context) int64 {
	userIDStr, isPass := common.ParseParam("user_id", client)
	if !isPass {
		return 0
	}
	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Fail(client, domain.ErrCodeJSONParseFailed, "参数解析失败")
		return 0
	}
	return userID
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
	if err != nil {
		var bizErr *domain.BusinessError
		if !errors.As(err, &bizErr) {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}
		response.Fail(client, bizErr.Code, bizErr.Msg)
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
	if err != nil {
		var bizErr *domain.BusinessError
		if errors.As(err, &bizErr) {
		} else {
			bizErr = &domain.BusinessError{
				Code: domain.ErrCodeUnknown,
				Msg:  "未知错误",
			}
		}
		response.Fail(client, bizErr.Code, bizErr.Msg)
		return
	}
	
	//成功响应
	response.Success(client, "用户信息查询成功", userInfo)
}