package domain

const (
	RoleUser     = 0 // 普通用户
	RoleMerchant = 1 // 商户
	RoleAdmin    = 2 // 管理员
)

type User struct {
	UserID       int64  `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Avatar       string `json:"avatar"` //头像
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Phone        string `json:"phone"`
	QQ           string `json:"qq"`
	Gender       string `json:"gender"` //性别
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Role         int    `json:"role" gorm:"default:0"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"` // 明文密码
}

type UserRegisterResponse struct {
	UserID int64 `json:"user_id"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRefreshATResponse struct {
	Token string `json:"token"`
}

type UserUpdatePasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type GetUserInfoResponse struct {
	UserID       int64  `json:"user_id;primaryKey;autoIncrement"`
	Avatar       string `json:"avatar"` //头像
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Phone        string `json:"phone"`
	QQ           string `json:"qq"`
	Gender       string `json:"gender"` //性别
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Username     string `json:"username"`
	Role         int    `json:"role"`
}
