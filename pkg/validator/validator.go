package validator

import (
	"errors"
	"fmt"
	"strings"
)

func TrimAndCheckEmpty(str, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(str)
	if trimmed == "" {
		return "", errors.New(fieldName + "不能为空")
	}
	return trimmed, nil
}

func CheckLengthRange(str string, fieldName string, min, max int) error {
	length := len(str)
	if length < min || length > max {
		return fmt.Errorf("%s长度需在%d-%d位之间（当前%d位）", fieldName, min, max, length)
	}
	return nil
}
