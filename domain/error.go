package domain

import "errors"

// 业务错误码
const (
	ErrCodeToBeNamed           = 9999  //待命名错误码
	ErrCodeSuccess             = 10000 // 成功
	ErrCodeParamInvalid        = 10001 // 参数错误（用户名重复、密码格式错误等）
	ErrCodeJSONParseFailed     = 10002 // JSON解析失败
	ErrCodeDBError             = 10003 // 数据库操作异常
	ErrCodeHashFailed          = 10004 // 哈希加密失败
	ErrCodeUnknown             = 10005 //未知错误
	ErrCodeSystemError         = 10006 //系统错误
	ErrCodeLoginFailed         = 10007 //登录错误
	ErrCodeUserNotFound        = 10008 //用户不存在
	ErrCodePasswordMismatch    = 10009 //密码不匹配
	ErrCodeTokenEmpty          = 10010 //token为空
	ErrCodeTokenFormatError    = 10011 //token格式错误
	ErrCodeTokenExpired        = 10012 //token过期
	ErrCodeTokenInvalid        = 10013 //token无效
	ErrCodeRefreshTokenEmpty   = 10014 //刷新AT时RT为空
	ErrCodeRefreshTokenExpired = 10015 //刷新AT时RT失效
	ErrCodePermissionDenied    = 10016 //权限不足
	ErrCodeMerchantExists      = 10017 //商户名已存在
	ErrCodeAlreadyMerchant     = 10018 //用户已是商户
	ErrCodeProductExists       = 10019 //商品名已存在
	ErrCodeProductNotFound     = 10020 //商品不存在
)

var (
	ErrTokenExpired     = errors.New("token_expired")
	ErrTokenInvalid     = errors.New("token_invalid")
	ErrTokenEmpty       = errors.New("token_empty")
	ErrTokenFormatError = errors.New("token_format_error")

	ErrAccessTokenClaims = errors.New("token_signature_invalid")

	ErrRefreshTokenEmpty   = errors.New("refresh_token_empty")
	ErrRefreshTokenExpired = errors.New("refresh_token_expired")
	ErrRefreshTokenClaims  = errors.New("refresh_token_signature_invalid")

	ErrTokenInvalidClaims  = errors.New("token_invalid_claims")
	ErrTokenInvalidExpType = errors.New("token_invalid_exp_type")
	ErrTokenSignFailed     = errors.New("token_sign_failed")
)

type BusinessError struct {
	Code int
	Msg  string
}

func (e *BusinessError) Error() string {
	return e.Msg
}
