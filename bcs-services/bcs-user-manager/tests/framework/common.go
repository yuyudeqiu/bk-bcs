package framework

import (
	"fmt"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	godbsdk "code.cwoa.net/carlchen2/cw-godb-sdk/core/config"
	gormsdk "code.cwoa.net/carlchen2/cw-godb-sdk/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

// TableNames 所有业务表的名称（按依赖顺序）
var TableNames = []string{
	"bcs_ws_cluster_credentials",
	"bcs_cluster_credentials",
	"bcs_register_tokens",
	"bcs_clusters",
	"bcs_client_users",
	"bcs_clients",
	"bcs_temp_tokens",
	"bcs_users",
	"bcs_user_resource_roles",
	"bcs_roles",
	"activities",
	"bcs_operation_logs",
	"bcs_token_notifies",
	"tke_cidrs",
}

// DatabaseSSL 数据库 SSL 配置
type DatabaseSSL struct {
	Enable      bool
	Mode        string
	Ca          string
	Cert        string
	Key         string
	KeyPassword string
}

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	DBType      string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	SvcConfPath string
	Ssl         DatabaseSSL
}

// NewDBClient 创建数据库客户端
func NewDBClient(cfg DatabaseConfig) (*gorm.DB, error) {
	var dbType godbsdk.DatabaseType
	switch cfg.DBType {
	case "mysql":
		dbType = godbsdk.Mysql
	case "gaussdb":
		dbType = godbsdk.Gaussdb
	case "dameng":
		dbType = godbsdk.Dameng
	default:
		dbType = godbsdk.Mysql
	}

	dbConfig := godbsdk.Database{
		Typex:    dbType,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		Ssl: godbsdk.TLS{
			Enable:      cfg.Ssl.Enable,
			Mode:        cfg.Ssl.Mode,
			Ca:          cfg.Ssl.Ca,
			Cert:        cfg.Ssl.Cert,
			Key:         cfg.Ssl.Key,
			KeyPassword: cfg.Ssl.KeyPassword,
		},
		SvcConfPath: cfg.SvcConfPath,
	}

	client, err := gormsdk.NewClient(dbConfig, logger.Default.LogMode(logger.Silent))
	if err != nil {
		return nil, fmt.Errorf("创建数据库客户端失败: %w", err)
	}
	return client.DB(), nil
}

// InitTables 使用 GORM AutoMigrate 创建所有表结构
// 支持 MySQL、达梦、OpenGauss
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.BcsUser{},
		&models.BcsTempToken{},
		&models.BcsClient{},
		&models.BcsCluster{},
		&models.BcsRegisterToken{},
		&models.BcsClusterCredential{},
		&models.BcsWsClusterCredentials{},
		&models.BcsRole{},
		&models.BcsUserResourceRole{},
		&models.Activity{},
		&models.BcsOperationLog{},
		&models.BcsTokenNotify{},
		&models.TkeCidr{},
	)
}

// CleanData 清理测试数据
func CleanData(db *gorm.DB) {
	for _, tableName := range TableNames {
		db.Exec("DELETE FROM " + tableName)
	}
}

// GetEnvWithFallback 获取环境变量，带有备用 key 的 fallback
func GetEnvWithFallback(key, fallbackKey string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return os.Getenv(fallbackKey)
}
