package framework

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// DMConfig 达梦测试配置
var DMConfig = DatabaseConfig{
	DBType:     "dameng",
	DBHost:     "localhost",
	DBPort:     5236,
	DBUser:     "bcs",
	DBPassword: "bcs123",
	DBName:     "BCS_TEST",
}

func init() {
	if host := GetEnvWithFallback("BCS_TEST_DM_HOST", "BKAUTH_TEST_DM_HOST"); host != "" {
		DMConfig.DBHost = host
	}
	if portStr := GetEnvWithFallback("BCS_TEST_DM_PORT", "BKAUTH_TEST_DM_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			DMConfig.DBPort = port
		}
	}
	if user := GetEnvWithFallback("BCS_TEST_DM_USER", "BKAUTH_TEST_DM_USER"); user != "" {
		DMConfig.DBUser = user
	}
	if pwd := GetEnvWithFallback("BCS_TEST_DM_PASSWORD", "BKAUTH_TEST_DM_PASSWORD"); pwd != "" {
		DMConfig.DBPassword = pwd
	}
	if name := GetEnvWithFallback("BCS_TEST_DM_DBNAME", "BKAUTH_TEST_DM_DBNAME"); name != "" {
		DMConfig.DBName = name
	}
	if tlsEnabledStr := GetEnvWithFallback("BCS_TEST_DM_TLS_ENABLED", "BKAUTH_TEST_DM_TLS_ENABLED"); tlsEnabledStr != "" {
		DMConfig.Ssl.Enable = (tlsEnabledStr == "true")
	}
	if certFile := GetEnvWithFallback("BCS_TEST_DM_TLS_CERT_FILE", "BKAUTH_TEST_DM_TLS_CERT_FILE"); certFile != "" {
		DMConfig.Ssl.Cert = certFile
	}
	if certKeyFile := GetEnvWithFallback("BCS_TEST_DM_TLS_CERT_KEY_FILE", "BKAUTH_TEST_DM_TLS_CERT_KEY_FILE"); certKeyFile != "" {
		DMConfig.Ssl.Key = certKeyFile
	}
	if svcConfPath := GetEnvWithFallback("BCS_TEST_DM_SVC_CONF_PATH", "BKAUTH_TEST_DM_SVC_CONF_PATH"); svcConfPath != "" {
		DMConfig.SvcConfPath = svcConfPath
	}
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
// dbName 为空时使用 DMConfig.DBName 的默认值
func MustInitDM(dbName string) *gorm.DB {
	if dbName == "" {
		dbName = DMConfig.DBName
	}
	db, err := InitDM(dbName)
	if err != nil {
		panic(fmt.Errorf("初始化达梦失败: %w", err))
	}
	return db
}
