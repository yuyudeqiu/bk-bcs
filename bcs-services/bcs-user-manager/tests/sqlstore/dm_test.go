//go:build dm

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

func TestTokenStoreDM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Store DM Integration Suite")
}

var db *gorm.DB

var _ = BeforeSuite(func() {
	db = framework.MustInitDM()

	err := framework.InitTables(db)
	Expect(err).ShouldNot(HaveOccurred())

	// 设置全局 DB
	sqlstore.SetGCoreDB(db)

	fmt.Println("Token Store DM 集成测试环境已就绪")
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
