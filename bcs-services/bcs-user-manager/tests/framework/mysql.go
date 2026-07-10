package framework

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var defaultMySQLTestDBPrefixes = []string{
	"bcs_user_test",
	"bcs_user_migration_test",
	"bcs_user_iam_migration_test",
}

// MySQLConfig MySQL 测试配置
var MySQLConfig = DatabaseConfig{
	DBType:     "mysql",
	DBHost:     "localhost",
	DBPort:     3306,
	DBUser:     "root",
	DBPassword: "root",
	DBName:     "bcs_test",
}

func init() {
	if host := GetEnvWithFallback("BCS_TEST_MYSQL_HOST", "BKAUTH_TEST_MYSQL_HOST"); host != "" {
		MySQLConfig.DBHost = host
	}
	if portStr := GetEnvWithFallback("BCS_TEST_MYSQL_PORT", "BKAUTH_TEST_MYSQL_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			MySQLConfig.DBPort = port
		}
	}
	if user := GetEnvWithFallback("BCS_TEST_MYSQL_USER", "BKAUTH_TEST_MYSQL_USER"); user != "" {
		MySQLConfig.DBUser = user
	}
	if pwd := GetEnvWithFallback("BCS_TEST_MYSQL_PASSWORD", "BKAUTH_TEST_MYSQL_PASSWORD"); pwd != "" {
		MySQLConfig.DBPassword = pwd
	}
	if name := GetEnvWithFallback("BCS_TEST_MYSQL_DBNAME", "BKAUTH_TEST_MYSQL_DBNAME"); name != "" {
		MySQLConfig.DBName = name
	}
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

// DropMySQLDatabaseIfExists 连接到系统数据库并删除目标数据库
func DropMySQLDatabaseIfExists(dbName string) error {
	if !isAllowedMySQLTestDBName(dbName) {
		return fmt.Errorf("refuse to drop database %q: set BCS_TEST_MYSQL_DBNAME_PREFIX to an allowed test prefix", dbName)
	}

	cfg := MySQLConfig
	cfg.DBName = "mysql" // 连接到 mysql 系统数据库

	db, err := NewDBClient(cfg)
	if err != nil {
		return fmt.Errorf("连接 mysql 系统数据库失败: %w", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error
	if err != nil {
		return fmt.Errorf("删除数据库 %s 失败: %w", dbName, err)
	}

	return nil
}

func isAllowedMySQLTestDBName(dbName string) bool {
	prefixes := append([]string{}, defaultMySQLTestDBPrefixes...)
	if configuredPrefixes := os.Getenv("BCS_TEST_MYSQL_DBNAME_PREFIX"); configuredPrefixes != "" {
		prefixes = append(prefixes, strings.Split(configuredPrefixes, ",")...)
	}

	for _, prefix := range prefixes {
		prefix = strings.TrimSpace(prefix)
		if prefix != "" && strings.HasPrefix(dbName, prefix) {
			return true
		}
	}
	return false
}

// MustInitMySQL 初始化 MySQL，失败时 panic
func MustInitMySQL(dbName string) *gorm.DB {
	db, err := InitMySQL(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化 MySQL 失败: %w", err))
	}
	return db
}
