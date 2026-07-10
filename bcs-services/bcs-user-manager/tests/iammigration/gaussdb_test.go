//go:build gaussdb

package iammigration_test

import (
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/tests/framework"
)

func TestIAMMigrateGaussDB(t *testing.T) {
	db := framework.MustInitGaussDB("")
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get raw sql db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	runIAMMigrationTest(t, db, "gaussdb")
}
