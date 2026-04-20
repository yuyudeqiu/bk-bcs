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
	DBName:     "bcs_test",
}

// InitGaussDB 初始化 OpenGauss 测试数据库
func InitGaussDB(dbName string) (*gorm.DB, error) {
	cfg := GaussDBConfig
	if dbName != "" {
		cfg.DBName = dbName
	}

	// 先尝试连接 postgres 默认数据库以创建目标数据库
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	db, err := NewDBClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 OpenGauss 失败: %w", err)
	}

	return db, nil
}

// createDatabaseIfNotExists 连接到默认数据库并创建目标数据库
func createDatabaseIfNotExists(cfg DatabaseConfig) error {
	adminCfg := cfg
	adminCfg.DBName = "postgres" // 连接到默认的 postgres 数据库

	db, err := NewDBClient(adminCfg)
	if err != nil {
		return fmt.Errorf("连接默认数据库 postgres 失败: %w", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 检查数据库是否存在
	var count int
	err = db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&count).Error
	if err != nil {
		return fmt.Errorf("检查数据库 %s 是否存在失败: %w", cfg.DBName, err)
	}

	if count == 0 {
		// 创建数据库
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)).Error
		if err != nil {
			return fmt.Errorf("创建数据库 %s 失败: %w", cfg.DBName, err)
		}
	}

	return nil
}

// MustInitGaussDB 初始化 OpenGauss，失败时 panic
func MustInitGaussDB(dbName string) *gorm.DB {
	db, err := InitGaussDB(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化 OpenGauss 失败: %w", err))
	}
	return db
}
