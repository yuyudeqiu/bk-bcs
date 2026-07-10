//go:build dm

package iammigration_test

import (
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/tests/framework"
)

func TestIAMMigrateDM(t *testing.T) {
	db := framework.MustInitDM("")
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get raw sql db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	runIAMMigrationTest(t, db, "dm")
}
