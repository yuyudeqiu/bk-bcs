//go:build !dm && !gaussdb

package sqlstore_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/tests/framework"
)

func TestTokenStoreMySQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Store MySQL Integration Suite")
}

var db *gorm.DB
var mysqlTestDBName string

var _ = BeforeSuite(func() {
	var err error
	mysqlTestDBName = framework.GetEnvWithFallback("BCS_TEST_MYSQL_DBNAME", "BKAUTH_TEST_MYSQL_DBNAME")
	if mysqlTestDBName == "" {
		mysqlTestDBName = "bcs_user_test"
	}

	db, err = framework.InitMySQL(mysqlTestDBName)
	Expect(err).ShouldNot(HaveOccurred())

	err = framework.InitTables(db)
	Expect(err).ShouldNot(HaveOccurred())

	// 设置全局 DB
	sqlstore.SetGCoreDB(db)

	fmt.Println("Token Store MySQL 集成测试环境已就绪")
})

var _ = AfterSuite(func() {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
	Expect(framework.DropMySQLDatabaseIfExists(mysqlTestDBName)).Should(Succeed())
})

var _ = BeforeEach(func() {
	framework.CleanData(db)
})

var _ = Describe("Token Store", func() {
	var s storeSet

	BeforeEach(func() {
		s.db = db
		s.tokenStore = sqlstore.NewTokenStore(db, nil)
		s.tokenNotifyStore = sqlstore.NewTokenNotifyStore(db)
	})

	describeStoreTests(&s)
})
