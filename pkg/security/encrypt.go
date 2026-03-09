package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword 密码加密
func HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
