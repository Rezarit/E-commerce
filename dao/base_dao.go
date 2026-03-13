package dao

import (
	"fmt"
	"log"
)

// InsertRecord 通用插入函数
func InsertRecord[T any](record *T) error {
	log.Printf("[DAO] 开始通用插入操作 | 类型：%T", record)
	if err := DB.Create(record).Error; err != nil {
		log.Printf("[DAO] 通用插入操作失败 | 类型：%T | 错误：%v", record, err)
		return err
	}
	log.Printf("[DAO] 通用插入操作成功 | 类型：%T", record)
	return nil
}

// UpdateRecord 通用更新函数
func UpdateRecord[T any](fieldName string, ID int64, record *T) error {
	log.Printf("[DAO] 开始通用更新操作 | ID：%d | 类型：%T", ID, record)
	if err := DB.Model(new(T)).
		Where(fieldName+" = ?", ID).
		Updates(record).Error; err != nil {
		log.Printf("[DAO] 通用更新操作失败 | ID：%d | 类型：%T | 错误：%v", ID, record, err)
		return err
	}
	log.Printf("[DAO] 通用更新操作成功 | ID：%d | 类型：%T", ID, record)
	return nil
}

// DeleteRecord 通用删除函数
func DeleteRecord[T any](fieldName string, ID int64) error {
	log.Printf("[DAO] 开始通用删除操作 | ID：%d | 类型：%T", ID, new(T))
	if err := DB.Model(new(T)).Where(fieldName+" = ?", ID).Delete(new(T)).Error; err != nil {
		log.Printf("[DAO] 通用删除操作失败 | ID：%d | 类型：%T | 错误：%v", ID, new(T), err)
		return err
	}
	log.Printf("[DAO] 通用删除操作成功 | ID：%d | 类型：%T", ID, new(T))
	return nil
}

// GetRecordByField 根据指定字段查询记录
func GetRecordByField[T any, V any](fieldName string, fieldValue V, result *T) error {
	log.Printf("[DAO] 开始根据字段查询记录 | 字段：%s | 值：%v | 结果类型：%T", fieldName, fieldValue, result)
	if err := DB.Model(new(T)).Where(fieldName+" = ?", fieldValue).Take(result).Error; err != nil {
		log.Printf("[DAO] 根据字段查询记录失败 | 字段：%s | 值：%v | 结果类型：%T | 错误：%v", fieldName, fieldValue, result, err)
		return err
	}
	log.Printf("[DAO] 根据字段查询记录成功 | 字段：%s | 值：%v | 结果类型：%T", fieldName, fieldValue, result)
	return nil
}

// GetRecordsByField 根据指定字段查询多条记录
func GetRecordsByField[T any, V any](fieldName string, fieldValue V, result *[]T) error {
	log.Printf("[DAO] 开始根据字段查询多条记录 | 字段：%s | 值：%v | 结果类型：%T", fieldName, fieldValue, result)
	if err := DB.Model(new(T)).Where(fieldName+" = ?", fieldValue).Find(result).Error; err != nil {
		log.Printf("[DAO] 根据字段查询多条记录失败 | 字段：%s | 值：%v | 结果类型：%T | 错误：%v", fieldName, fieldValue, result, err)
		return err
	}
	log.Printf("[DAO] 根据字段查询多条记录成功 | 字段：%s | 值：%v | 结果类型：%T", fieldName, fieldValue, result)
	return nil
}

// CheckFieldExists 检查指定字段的值是否存在
func CheckFieldExists[T any, V any](fieldName string, fieldValue V) (bool, error) {
	log.Printf("[DAO] 开始检查字段值是否存在 | 字段：%s | 值：%v | 类型：%T", fieldName, fieldValue, new(T))
	var count int64
	result := DB.Model(new(T)).Where(fieldName+" = ?", fieldValue).Count(&count)

	if result.Error != nil {
		log.Printf("[DAO] 检查字段值存在性失败 | 字段：%s | 值：%v | 类型：%T | 错误：%v", fieldName, fieldValue, new(T), result.Error)
		return false, fmt.Errorf("检查字段值存在性失败: %w", result.Error)
	}

	exists := count > 0
	log.Printf("[DAO] 检查字段值存在性成功 | 字段：%s | 值：%v | 类型：%T | 存在：%v", fieldName, fieldValue, new(T), exists)
	return exists, nil
}
