package domain

type User struct {
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
	Password     string `json:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"` // 明文密码
}

type UserRegisterResponse struct {
	Username string `json:"username"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
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
}
