package dao

import (
	"fmt"
	"github.com/Rezarit/E-commerce/domain"
)

// CheckUsernameExists 用户名查询
func CheckUsernameExists(username string) (bool, error) {
	var count int64
	result := DB.Model(&domain.User{}).Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("查询用户名失败: %w", result.Error)
	}
	return count > 0, nil
}

// InsertRecord 通用插入函数
func InsertRecord[T any](record *T) error {
	return DB.Create(record).Error
}

// InsertUser 用户相关数据插入数据库
func InsertUser(user *domain.User) error {
	return InsertRecord(user)
}

// UpdateRecord 通用更新函数
func UpdateRecord[T any](ID int64, record *T) error {
	return DB.Model(new(T)).
		Where("user_id = ?", ID).
		Updates(record).Error
}

// UpdateUser 更新用户信息
func UpdateUser(userID int64, user *domain.User) error {
	return UpdateRecord(userID, user)
}

// GetUserPasswordByUsername 通过用户名查密码
func GetUserPasswordByUsername(username string) (string, error) {
	var user domain.User
	err := DB.Model(&domain.User{}).Select("password").Where("username = ?", username).Take(&user).Error
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

// GetUserIDByUsername 通过用户名查询用户ID
func GetUserIDByUsername(username string) (int64, error) {
	var user domain.User
	err := DB.Model(&domain.User{}).Select("user_id").Where("username = ?", username).Take(&user).Error
	if err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// GetUserInfoByID 通过用户ID查询用户信息
func GetUserInfoByID(userID int64) (domain.GetUserInfoResponse, error) {
	var user domain.User
	err := DB.Model(&domain.User{}).Where("user_id = ?", userID).Take(&user).Error
	if err != nil {
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
	}, nil
}

//func UpdatePassword(id int, password string) error {
//	cmd := "UPDATE users SET password=? WHERE Id=?;"
//	_, err := DB.Exec(cmd, password, id)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func SearchUserMsg(ID int) (domain.User, error) {
//	var user domain.User
//	cmd := "SELECT * FROM users WHERE Id=?;"
//	err := DB.QueryRow(cmd, ID).Scan(
//		&user.UserID,
//		&user.Avatar,
//		&user.Nickname,
//		&user.Introduction,
//		&user.Phone,
//		&user.QQ,
//		&user.Gender,
//		&user.Email,
//		&user.Birthday)
//	if err != nil {
//		return domain.User{}, err
//	}
//	return user, nil
//}
//
//func UpdateUserMeg(UserID, Phone, QQ int, Avatar, Nickname, Introduction, Gender, Email, Birthday string) error {
//	cmd := "UPDATE users SET avatar = ?, nickname = ?, introduction = ?, phone = ?, qq = ?, gender = ?, email = ?, birthday = ? WHERE id = ?"
//	_, err := DB.Exec(cmd, Avatar, Nickname, Introduction, Phone, QQ, Gender, Email, Birthday, UserID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
