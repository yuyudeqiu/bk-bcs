package app

import (
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/options"
)

func TestParseDatabaseConnectionConfig(t *testing.T) {
	t.Run("structured database config without legacy DSN", func(t *testing.T) {
		op := &options.UserManagerOptions{DatabaseConfig: options.DatabaseConfig{
			DBType:     "goldendb",
			DBHost:     "127.0.0.1",
			DBPort:     3306,
			DBUser:     "root",
			DBPassword: "password",
			DBName:     "bcs",
		}}
		conf := config.NewUserMgrConfig()

		if err := parseDatabaseConnectionConfig(op, conf); err != nil {
			t.Fatalf("parse structured database config: %v", err)
		}
		if conf.DSN != "" {
			t.Fatalf("DSN = %q, want empty", conf.DSN)
		}
		if conf.DatabaseConfig.DBType != "goldendb" {
			t.Fatalf("database type = %q, want goldendb", conf.DatabaseConfig.DBType)
		}
	})

	t.Run("legacy DSN takes precedence", func(t *testing.T) {
		op := &options.UserManagerOptions{
			DSN: "root:password@tcp(127.0.0.1:3306)/bcs",
			DatabaseConfig: options.DatabaseConfig{
				DBType: "goldendb",
			},
		}
		conf := config.NewUserMgrConfig()

		if err := parseDatabaseConnectionConfig(op, conf); err != nil {
			t.Fatalf("parse legacy DSN: %v", err)
		}
		if conf.DSN != op.DSN {
			t.Fatalf("DSN = %q, want %q", conf.DSN, op.DSN)
		}
		if conf.DatabaseConfig.DBType != "" {
			t.Fatalf("structured database config should not be populated when DSN is configured")
		}
	})
}
