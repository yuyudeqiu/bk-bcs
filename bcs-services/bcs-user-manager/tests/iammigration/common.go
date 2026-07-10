package iammigration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tencent/bk-bcs/bcs-common/pkg/auth/iam"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	usermanager "github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/migrations"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/options"
)

const iamMigrationTestTable = "bk_iam_migrations"

// runIAMMigrationTest is shared by all supported database test entry points.
func runIAMMigrationTest(t *testing.T, db *gorm.DB, dbType string) {
	t.Helper()

	if err := db.Migrator().DropTable(iamMigrationTestTable); err != nil {
		t.Fatalf("drop previous migration table: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Migrator().DropTable(iamMigrationTestTable); err != nil {
			t.Errorf("drop migration table: %v", err)
		}
	})

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get raw sql db: %v", err)
	}

	iamServer := newIAMMigrationTestServer(t)
	t.Cleanup(iamServer.Close)

	conf := &config.UserMgrConfig{
		ServCert:       &config.CertConfig{},
		DatabaseConfig: options.DatabaseConfig{DBType: dbType},
		IAMConfig: options.IAMConfig{
			SystemID:    "bk_bcs",
			AppCode:     "bk_bcs",
			AppSecret:   "mock_secret",
			GateWayHost: iamServer.URL,
		},
	}
	permClient, err := iam.NewIamMigrateClient(&iam.Options{
		SystemID:    conf.IAMConfig.SystemID,
		AppCode:     conf.IAMConfig.AppCode,
		AppSecret:   conf.IAMConfig.AppSecret,
		GateWayHost: conf.IAMConfig.GateWayHost,
	})
	if err != nil {
		t.Fatalf("create iam migrate client: %v", err)
	}

	manager := usermanager.NewUserManager(conf)
	manager.IamPermClient = permClient

	sourceDriver, err := iofs.New(migrations.MigrationFS, ".")
	if err != nil {
		t.Fatalf("open migration source: %v", err)
	}

	templateVar := map[string]string{
		"BK_IAM_SYSTEM_ID": conf.IAMConfig.SystemID,
		"APP_CODE":         conf.IAMConfig.AppCode,
		"BCS_HOST":         "http://bcs.example.com",
	}
	if err := manager.MigrateIAM(sqlDB, sourceDriver, iamMigrationTestTable, 5*time.Minute, templateVar); err != nil {
		t.Fatalf("run iam migration through UserManager: %v", err)
	}

	var version int
	var dirty bool
	if err := db.Table(iamMigrationTestTable).Select("version, dirty").Row().Scan(&version, &dirty); err != nil {
		t.Fatalf("query migration version: %v", err)
	}
	if dirty {
		t.Fatalf("migration table is dirty at version %d", version)
	}
	if version != 13 {
		t.Fatalf("unexpected migration version: got %d, want 13", version)
	}
}

func newIAMMigrationTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data":    map[string]interface{}{},
		}); err != nil {
			t.Errorf("write fake iam response: %v", err)
		}
	}))
}
