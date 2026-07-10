//go:build !dm && !gaussdb

package iammigration_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/tests/framework"
)

func TestIAMMigrateMySQLCompatible(t *testing.T) {
	dbType := strings.ToLower(framework.GetEnvWithFallback("BCS_TEST_DB_TYPE", "BKAUTH_TEST_DB_TYPE"))
	if dbType == "" {
		dbType = "mysql"
	}
	if dbType != "mysql" && dbType != "ob" && dbType != "oceanbase" {
		t.Fatalf("unsupported MySQL-compatible database type %q", dbType)
	}

	dbName := framework.GetEnvWithFallback("BCS_TEST_MYSQL_DBNAME", "BKAUTH_TEST_MYSQL_DBNAME")
	if dbName == "" {
		dbName = "bcs_user_iam_migration_test"
	}
	db, err := framework.InitMySQL(dbName)
	if err != nil {
		t.Fatalf("init mysql-compatible database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get raw sql db: %v", err)
	}
	t.Cleanup(func() { closeAndDropMySQL(t, sqlDB, dbName) })

	runIAMMigrationTest(t, db, dbType)
}

func closeAndDropMySQL(t *testing.T, sqlDB *sql.DB, dbName string) {
	t.Helper()
	if sqlDB != nil {
		_ = sqlDB.Close()
	}
	if err := framework.DropMySQLDatabaseIfExists(dbName); err != nil {
		t.Errorf("drop mysql database %s: %v", dbName, err)
	}
}
