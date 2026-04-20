package framework

import (
	"fmt"

	"gorm.io/gorm"
)

// GaussDBConfig OpenGauss 测试配置
var GaussDBConfig = DatabaseConfig{
	DBType:     "gaussdb",
	DBHost:     "localhost",
	DBPort:     5433,
	DBUser:     "gaussdb",
	DBPassword: "openGauss@123",
	DBName:     "gorm_test",
}

// InitGaussDB 初始化 OpenGauss 测试数据库
func InitGaussDB(dbName string) (*gorm.DB, error) {
	cfg := GaussDBConfig
	if dbName != "" {
		cfg.DBName = dbName
	}

	db, err := NewDBClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 OpenGauss 失败: %w", err)
	}

	return db, nil
}

// MustInitGaussDB 初始化 OpenGauss，失败时 panic
func MustInitGaussDB(dbName string) *gorm.DB {
	db, err := InitGaussDB(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化 OpenGauss 失败: %w", err))
	}
	return db
}
