package dao

import (
	"github.com/Rezarit/go-seckill-system/domain"
)

// CheckUsernameExists 用户名是否存在
func CheckUsernameExists(username string) (bool, error) {
	exists, err := CheckFieldExists[domain.User, string]("username", username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertUser 用户相关数据插入数据库
func InsertUser(user *domain.User) error {
	if err := InsertRecord(user); err != nil {
		return err
	}
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(userID int64, user *domain.User) error {
	if err := UpdateRecord("user_id", userID, user); err != nil {
		return err
	}
	return nil
}

// GetUserPasswordByUsername 通过用户名查密码
func GetUserPasswordByUsername(username string) (string, error) {
	var user domain.User
	err := GetRecordByField[domain.User, string]("username", username, &user)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

// GetUserIDByUsername 通过用户名查询用户ID
func GetUserIDByUsername(username string) (int64, error) {
	var user domain.User
	err := GetRecordByField[domain.User, string]("username", username, &user)
	if err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// GetUserInfoByID 通过用户ID查询用户信息
func GetUserInfoByID(userID int64) (domain.GetUserInfoResponse, error) {
	var user domain.User
	if err := GetRecordByField("user_id", userID, &user); err != nil {
		return domain.GetUserInfoResponse{}, err
	}
	return domain.GetUserInfoResponse{
		UserID:       user.UserID,
		Avatar:       user.Avatar,
		Nickname:     user.Nickname,
		Introduction: user.Introduction,
		Phone:        user.Phone,
		QQ:           user.QQ,
		Gender:       user.Gender,
		Email:        user.Email,
		Birthday:     user.Birthday,
		Username:     user.Username,
		Role:         user.Role,
	}, nil
}
