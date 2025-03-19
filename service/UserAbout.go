package service

import "golang.org/x/crypto/bcrypt"

// CheckUsernameLen 检查用户名长度
func CheckUsernameLen(username string) int {
	if 1 > len(username) {
		return 10001
	}
	return 10000
}

// CheckPasswordLen 检查密码长度
func CheckPasswordLen(password string) int {
	if 8 > len(password) {
		return 10001
	}
	return 10000
}

// HashedPassword 密码加密
func HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
