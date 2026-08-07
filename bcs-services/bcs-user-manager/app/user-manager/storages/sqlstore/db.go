/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sqlstore

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	godbsdk "code.cwoa.net/carlchen2/cw-godb-sdk/core/config"
	gormsdk "code.cwoa.net/carlchen2/cw-godb-sdk/gorm"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/config"
)

// GCoreDB global DB client
var GCoreDB *gorm.DB

// SetGCoreDB 设置全局 DB（用于测试）
func SetGCoreDB(db *gorm.DB) {
	GCoreDB = db
}

// InitCoreDatabase set DB client
func InitCoreDatabase(conf *config.UserMgrConfig) error {
	if conf == nil {
		return fmt.Errorf("core_database config not init")
	}

	var db *gorm.DB
	var err error

	// 优先使用 DSN 方式（向后兼容）
	if conf.DSN != "" {
		db, err = gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.SetConnMaxLifetime(60 * time.Second)
		sqlDB.SetMaxIdleConns(20)
		sqlDB.SetMaxOpenConns(20)
	} else if conf.DatabaseConfig.DBHost != "" {
		// 使用 SDK 结构化配置
		dbConfig := godbsdk.Database{
			Typex:    sdkDatabaseType(conf.DatabaseConfig.DBType),
			Host:     conf.DatabaseConfig.DBHost,
			Port:     conf.DatabaseConfig.DBPort,
			User:     conf.DatabaseConfig.DBUser,
			Password: conf.DatabaseConfig.DBPassword,
			Name:     conf.DatabaseConfig.DBName,
			Ssl: godbsdk.TLS{
				Enable:      conf.DatabaseConfig.Ssl.Enable,
				Mode:        conf.DatabaseConfig.Ssl.Mode,
				Ca:          conf.DatabaseConfig.Ssl.Ca,
				Cert:        conf.DatabaseConfig.Ssl.Cert,
				Key:         conf.DatabaseConfig.Ssl.Key,
				KeyPassword: conf.DatabaseConfig.Ssl.KeyPassword,
			},
			SvcConfPath: conf.DatabaseConfig.SvcConfPath,
		}
		// 连接池参数使用配置值，若为 0 则使用 SDK 默认值
		if conf.DatabaseConfig.MaxOpenConns > 0 {
			dbConfig.MaxOpenConns = conf.DatabaseConfig.MaxOpenConns
		}
		if conf.DatabaseConfig.MaxIdleConns > 0 {
			dbConfig.MaxIdleConns = conf.DatabaseConfig.MaxIdleConns
		}
		if conf.DatabaseConfig.ConnMaxLifetimeSecond > 0 {
			dbConfig.ConnMaxLifetimeSecond = conf.DatabaseConfig.ConnMaxLifetimeSecond
		}

		client, err := gormsdk.NewClient(dbConfig, logger.Default.LogMode(logger.Silent))
		if err != nil {
			return err
		}
		db = client.DB()
	} else {
		return fmt.Errorf("core_database dsn not configured and database_config is empty")
	}

	if db == nil {
		return fmt.Errorf("database initialized with nil gorm DB")
	}

	GCoreDB = db
	return nil
}

func sdkDatabaseType(dbType string) godbsdk.DatabaseType {
	switch strings.ToLower(dbType) {
	case "mysql", "ob", "oceanbase", "dg", "goldendb":
		return godbsdk.Mysql
	case "postgres":
		return godbsdk.Postgres
	case "dameng":
		return godbsdk.Dameng
	default:
		return godbsdk.Mysql
	}
}
