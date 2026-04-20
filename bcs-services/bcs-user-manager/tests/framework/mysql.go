package framework

import (
	"fmt"

	"gorm.io/gorm"
)

// MySQLConfig MySQL 测试配置
var MySQLConfig = DatabaseConfig{
	DBType:    "mysql",
	DBHost:     "localhost",
	DBPort:     3306,
	DBUser:     "root",
	DBPassword: "root",
	DBName:     "bcs_test",
}

// InitMySQL 初始化 MySQL 测试数据库
func InitMySQL(dbName string) (*gorm.DB, error) {
	cfg := MySQLConfig
	if dbName != "" {
		cfg.DBName = dbName
	}

	// 先尝试连接 mysql 默认数据库以创建目标数据库
	if err := createMySQLDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	db, err := NewDBClient(cfg)
	if err != nil {
		return nil, err
	}

	// 设置 sql_mode
	sqlDB, _ := db.DB()
	_, _ = sqlDB.Exec("SET sql_mode=''")

	return db, nil
}

// createMySQLDatabaseIfNotExists 连接到系统数据库并创建目标数据库
func createMySQLDatabaseIfNotExists(cfg DatabaseConfig) error {
	adminCfg := cfg
	adminCfg.DBName = "mysql" // 连接到 mysql 系统数据库

	db, err := NewDBClient(adminCfg)
	if err != nil {
		return fmt.Errorf("连接 mysql 系统数据库失败: %w", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 创建数据库 (MySQL 支持 IF NOT EXISTS)
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName)).Error
	if err != nil {
		return fmt.Errorf("创建数据库 %s 失败: %w", cfg.DBName, err)
	}

	return nil
}

// MustInitMySQL 初始化 MySQL，失败时 panic
func MustInitMySQL(dbName string) *gorm.DB {
	db, err := InitMySQL(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化 MySQL 失败: %w", err))
	}
	return db
}
