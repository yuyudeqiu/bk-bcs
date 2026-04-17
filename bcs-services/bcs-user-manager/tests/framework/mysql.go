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
	DBName:     "bcs_user_test",
}

// InitMySQL 初始化 MySQL 测试数据库
func InitMySQL(dbName string) (*gorm.DB, error) {
	cfg := MySQLConfig
	cfg.DBName = dbName

	db, err := NewDBClient(cfg)
	if err != nil {
		return nil, err
	}

	// 设置 sql_mode
	sqlDB, _ := db.DB()
	_, _ = sqlDB.Exec("SET sql_mode=''")

	return db, nil
}

// MustInitMySQL 初始化 MySQL，失败时 panic
func MustInitMySQL(dbName string) *gorm.DB {
	db, err := InitMySQL(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化 MySQL 失败: %w", err))
	}
	return db
}
