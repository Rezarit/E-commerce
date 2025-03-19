package dao

import (
	"github.com/Rezarit/E-commerce/domain"
)

func CheckUsernameExists(username string) (bool, error) {
	var count int
	cmd := "SELECT COUNT(*) FROM users WHERE username =?"
	err := DB.QueryRow(cmd, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertPassword 执行插入指令
func InsertPassword(nickname, username, password string) error {
	cmd := "INSERT INTO users(nickname,username,password) VALUES (?,?,?,?);"
	_, err := DB.Exec(cmd, nickname, username, password)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePassword(id int, password string) error {
	cmd := "UPDATE users SET password=? WHERE Id=?;"
	_, err := DB.Exec(cmd, password, id)
	if err != nil {
		return err
	}
	return nil
}

func SearchUserMsg(ID int) (domain.User, error) {
	var user domain.User
	cmd := "SELECT * FROM users WHERE Id=?;"
	err := DB.QueryRow(cmd, ID).Scan(
		&user.UserID,
		&user.Avatar,
		&user.Nickname,
		&user.Introduction,
		&user.Phone,
		&user.QQ,
		&user.Gender,
		&user.Email,
		&user.Birthday)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func UpdateUserMeg(UserID, Phone, QQ int, Avatar, Nickname, Introduction, Gender, Email, Birthday string) error {
	cmd := "UPDATE users SET avatar = ?, nickname = ?, introduction = ?, phone = ?, qq = ?, gender = ?, email = ?, birthday = ? WHERE id = ?"
	_, err := DB.Exec(cmd, Avatar, Nickname, Introduction, Phone, QQ, Gender, Email, Birthday, UserID)
	if err != nil {
		return err
	}
	return nil
}