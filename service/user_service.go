package service

import (
	"errors"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
	"github.com/Rezarit/E-commerce/pkg/security"
	"github.com/Rezarit/E-commerce/pkg/token"
	"github.com/Rezarit/E-commerce/pkg/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func Register(req domain.UserRegisterRequest) (*domain.User, error) {
	//校验用户名
	if err := CheckUsername(req.Username); err != nil {
		return nil, err
	}
	// 校验密码
	if err := CheckPassword(req.Password); err != nil {
		return nil, err
	}
	// 校验用户名是否存在
	if err := CheckUsernameExists(req.Username); err != nil {
		return nil, err
	}

	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeHashFailed,
			Msg:  "密码加密失败，请稍后再试",
		}
	}

	// 整合信息
	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPwd),
	}

	// 执行插入指令
	if err = InsertUser(user); err != nil {
		return nil, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "用户注册失败，请稍后再试",
		}
	}

	return user, nil
}

// InsertUser 向数据库插入用户信息
func InsertUser(user *domain.User) error {
	if err := dao.InsertUser(user); err != nil {
		log.Printf("[Service] 用户数据入库失败 | 用户名：%s | 错误：%v", user.Username, err)
		return err
	}
	log.Printf("[Service] 用户数据入库成功 | 用户名：%s", user.Username)
	return nil
}

// Login 登录
func Login(req domain.UserLoginRequest) (userID int64, accessToken, refreshToken string, err error) {
	// 校验用户名
	if err = CheckUsername(req.Username); err != nil {
		log.Printf("[Service] 用户名参数校验失败 | 用户名：%s | 错误：%v", req.Username, err)
		return 0, "", "", err
	}
	// 校验密码
	if err = CheckPassword(req.Password); err != nil {
		log.Printf("[Service] 密码参数校验失败 | 用户名：%s | 错误：%v", req.Username, err)
		return 0, "", "", err
	}

	// 获取数据库加密密码
	hashedPassword, err := GetPasswordByUsername(req.Username)
	if err != nil {
		log.Printf("[Service] 数据库密码获取失败 | 用户名：%s | 错误：%v", req.Username, err)
		return 0, "", "", err
	}
	// 核对密码
	err = ComparePassword(hashedPassword, req.Password)
	if err != nil {
		log.Printf("[Service] 密码核对失败 | 用户名：%s | 错误：%v", req.Username, err)
		return 0, "", "", err
	}

	// 获取用户ID
	userID, err = GetUserIDByUsername(req.Username)
	if err != nil {
		log.Printf("[Service] 用户ID查询失败 | 用户名：%s | 错误：%v", req.Username, err)
		return 0, "", "", err
	}

	// 生成token和refresh_token
	accessToken, err = GenerateAccessToken(userID)
	if err != nil {
		log.Printf("[Service] AccessToken生成失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, "", "", err
	}
	refreshToken, err = GenerateRefreshToken(userID)
	if err != nil {
		log.Printf("[Service] RefreshToken生成失败 | 用户ID：%d | 错误：%v", userID, err)
		return 0, "", "", err
	}

	//构造响应结构体
	return userID, accessToken, refreshToken, nil
}

// GetPasswordByUsername 通过用户名查询密码
func GetPasswordByUsername(username string) (string, error) {
	password, err := dao.GetUserPasswordByUsername(username)
	if err != nil {
		// 仅处理“没找到用户”的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[Service] 密码查询失败 | 用户名：%s | 错误：用户不存在", username)
			return "", &domain.BusinessError{
				Code: domain.ErrCodeUserNotFound,
				Msg:  "用户名或密码错误",
			}
		}
		// 处理其他数据库错误
		log.Printf("[Service] 密码查询失败 | 用户名：%s | 错误：%v", username, err)
		return "", &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "密码查询失败，请稍后再试",
		}
	}
	return password, nil
}

func GetUserIDByUsername(username string) (int64, error) {
	userID, err := dao.GetUserIDByUsername(username)
	if err != nil {
		// 仅处理“没找到用户”的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[Service] 用户ID查询失败 | 用户名：%s | 错误：用户不存在", username)
			return 0, &domain.BusinessError{
				Code: domain.ErrCodeUserNotFound,
				Msg:  "用户名或密码错误",
			}
		}
		// 处理其他数据库错误
		log.Printf("[Service] 用户ID查询失败 | 用户名：%s | 错误：%v", username, err)
		return 0, &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "ID查询失败，请稍后再试",
		}
	}
	return userID, nil
}

// ComparePassword 核对密码
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("[Service] 密码核对失败 | 错误：密码与哈希值不匹配")
			return &domain.BusinessError{
				Code: domain.ErrCodePasswordMismatch,
				Msg:  "用户名或密码错误",
			}
		}

		// 处理其他错误
		log.Printf("[Service] 密码校验系统异常 | 错误：%v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeSystemError,
			Msg:  "密码校验失败，请联系客服",
		}
	}
	return nil
}

// CheckUsername 检查用户名是否符合要求
func CheckUsername(username string) error {
	trimmedName, err := validator.TrimAndCheckEmpty(username, "用户名")
	if err != nil {
		log.Printf("[Service] 用户名参数校验失败 | 用户名：%s | 错误：%v", username, err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if err = validator.CheckLengthRange(trimmedName, "用户名", 6, 20); err != nil {
		log.Printf("[Service] 用户名长度校验失败 | 用户名：%s | 错误：%v", username, err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	return nil
}

// CheckPassword 检查密码长度
func CheckPassword(password string) error {
	trimmedPassword, err := validator.TrimAndCheckEmpty(password, "密码")
	if err != nil {
		log.Printf("[Service] 密码参数校验失败 | 错误：%v", err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	if err = validator.CheckLengthRange(trimmedPassword, "密码", 8, 20); err != nil {
		log.Printf("[Service] 密码长度校验失败 | 错误：%v", err)
		return &domain.BusinessError{Code: domain.ErrCodeParamInvalid, Msg: err.Error()}
	}
	return nil
}

// CheckUsernameExists 用户名查重
func CheckUsernameExists(username string) error {
	isExists, err := dao.CheckUsernameExists(username)
	if err != nil {
		log.Printf("[Service] 用户名查重失败 | 用户名：%s | 错误：%v", username, err)
		return &domain.BusinessError{
			Code: domain.ErrCodeDBError,
			Msg:  "服务器出现问题，请稍后再试",
		}
	}
	if isExists {
		log.Printf("[Service] 用户名已存在 | 用户名：%s", username)
		return &domain.BusinessError{
			Code: domain.ErrCodeParamInvalid,
			Msg:  "该用户名已被注册，请更换",
		}
	}
	return nil
}

// HashPassword 密码加密
func HashPassword(password string) ([]byte, error) {
	hashedPwd, err := security.HashedPassword(password)
	if err != nil {
		log.Printf("[Service] 密码加密失败 | 错误：%v", err)
		return nil, err
	}
	return hashedPwd, nil
}

// CheckAuthHeader 检查authHeader是否为空
func CheckAuthHeader(authHeader string) error {
	_, err := validator.TrimAndCheckEmpty(authHeader, "token")
	if err != nil {
		log.Printf("[Service] Token刷新失败 | 错误：传入token为空")
		return &domain.BusinessError{
			Code: domain.ErrCodeRefreshTokenEmpty,
			Msg:  "传入token为空",
		}
	}
	return nil
}

func GetTokenFromAuthHeader(authHeader string) (string, error) {
	tokenString, err := token.GetTokenFromAuthHeader(authHeader)
	if err != nil {
		log.Printf("[Service] Token解析失败 | 错误：%v", err)
		return "", &domain.BusinessError{
			Code: domain.ErrCodeTokenFormatError,
			Msg:  err.Error(),
		}
	}
	return tokenString, nil
}

func ValidateRefreshToken(claims *domain.RefreshTokenClaims) error {
	err := token.ValidateRefreshToken(claims)
	if err != nil {
		log.Printf("[Service] RefreshToken校验失败 | 错误：%v", err)
		return &domain.BusinessError{
			Code: domain.ErrCodeRefreshTokenExpired,
			Msg:  err.Error(),
		}
	}
	return nil
}

// RefreshAccessToken 刷新AT
func RefreshAccessToken(authHeader string) (string, error) {
	err := CheckAuthHeader(authHeader)
	if err != nil {
		log.Printf("[Service] AT刷新参数校验失败 | 错误：%v", err)
		return "", err
	}

	tokenString, err := GetTokenFromAuthHeader(authHeader)
	if err != nil {
		log.Printf("[Service] AT刷新RT获取失败 | 错误：%v", err)
		return "", err
	}

	parsedToken, err := ParseRefreshToken(tokenString)
	if err != nil {
		log.Printf("[Service] AT刷新RT解析失败 | 错误：%v", err)
		return "", err
	}

	err = ValidateRefreshToken(parsedToken)
	if err != nil {
		log.Printf("[Service] AT刷新RT校验失败 | 用户ID：%d | 错误：%v", parsedToken.UserID, err)
		return "", err
	}

	accessToken, err := GenerateAccessToken(parsedToken.UserID)
	if err != nil {
		log.Printf("[Service] AT刷新AT生成失败 | 用户ID：%d | 错误：%v", parsedToken.UserID, err)
		return "", err
	}
	return accessToken, nil
}

func UpdateUserPassword(userID int64, username, oldPassword, newPassword string) error {
	// 从数据库中获取用户密码
	hashedPassword, err := dao.GetUserPasswordByUsername(username)
	if err != nil {
		log.Printf("[Service] 更新密码-密码获取失败 | 用户名：%s | 错误：%v", username, err)
		return err
	}

	// 核对旧密码
	err = ComparePassword(hashedPassword, oldPassword)
	if err != nil {
		log.Printf("[Service] 更新密码-旧密码校验失败 | 用户名：%s | 错误：%v", username, err)
		return err
	}

	// 加密新密码
	hashedNewPassword, err := HashPassword(newPassword)
	if err != nil {
		log.Printf("[Service] 更新密码-新密码加密失败 | 用户名：%s | 错误：%v", username, err)
		return err
	}

	// 更新数据库密码
	err = dao.UpdateUser(userID, &domain.User{
		Password: string(hashedNewPassword),
	})
	if err != nil {
		log.Printf("[Service] 更新密码新密码入库失败 | 用户名：%s | 错误：%v", username, err)
		return err
	}

	return nil
}

func UpdateUserInfoByID(userID int64, user domain.User) error {
	err := dao.UpdateUser(userID, &user)
	if err != nil {
		log.Printf("[Service] 更新用户信息失败 | 用户ID：%d | 错误：%v", userID, err)
		return err
	}
	return nil
}

func GetUserInfoById(userID int64) (domain.GetUserInfoResponse, error) {
	//从数据库中获取用户信息
	userInfo, err := dao.GetUserInfoByID(userID)
	if err != nil {
		log.Printf("[Service] 获取用户信息失败 | 用户ID：%d | 错误：%v", userID, err)
		return domain.GetUserInfoResponse{}, err
	}
	return userInfo, nil
}
