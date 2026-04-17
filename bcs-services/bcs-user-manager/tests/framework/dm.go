package framework

import (
	"fmt"

	"gorm.io/gorm"
)

// DMConfig 达梦测试配置
var DMConfig = DatabaseConfig{
	DBType:    "dameng",
	DBHost:     "localhost",
	DBPort:     5236,
	DBUser:     "bkauth",
	DBPassword: "bkauth123",
	DBName:     "BKAUTH",
}

// InitDM 初始化达梦测试数据库
func InitDM(dbName string) (*gorm.DB, error) {
	cfg := DMConfig
	cfg.DBName = dbName

	db, err := NewDBClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化达梦失败: %w", err)
	}

	return db, nil
}

// MustInitDM 初始化达梦，失败时 panic
func MustInitDM(dbName string) *gorm.DB {
	db, err := InitDM(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化达梦失败: %w", err))
	}
	return db
}
