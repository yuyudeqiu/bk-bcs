//go:build gaussdb

package sqlstore_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/tests/framework"
)

func TestTokenStoreGaussDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Store GaussDB Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error
	db, err = framework.InitGaussDB("bcs_user_test")
	Expect(err).ShouldNot(HaveOccurred())

	err = framework.InitTables(db)
	Expect(err).ShouldNot(HaveOccurred())

	fmt.Println("Token Store GaussDB 集成测试环境已就绪")
})

var _ = BeforeEach(func() {
	framework.CleanData(db)
	sqlstore.SetGCoreDB(db)
})

var (
	db         *gorm.DB
	tokenStore sqlstore.TokenStore
)

// CreateTestToken 创建测试 Token
func CreateTestToken(name string, userType uint, expired bool) *models.BcsUser {
	expiresAt := time.Now().Add(time.Hour)
	if expired {
		expiresAt = time.Now().Add(-time.Hour)
	}
	return &models.BcsUser{
		Name:      name,
		UserType:  userType,
		UserToken: "test-token-" + name,
		CreatedBy: "system",
		ExpiresAt: expiresAt,
	}
}

var _ = Describe("Token Store GaussDB 集成测试", func() {
	BeforeEach(func() {
		tokenStore = sqlstore.NewTokenStore(db, nil)
	})

	Describe("GetTokenByCondition", func() {
		It("按名称查询 Token", func() {
			token := CreateTestToken("test-user", models.PlainUser, false)
			err := tokenStore.CreateToken(token)
			Expect(err).ShouldNot(HaveOccurred())

			found := tokenStore.GetTokenByCondition(&models.BcsUser{Name: "test-user"})
			Expect(found).ShouldNot(BeNil())
			Expect(found.Name).Should(Equal("test-user"))
		})

		It("查询不存在的 Token 返回 nil", func() {
			found := tokenStore.GetTokenByCondition(&models.BcsUser{Name: "non-existent"})
			Expect(found).Should(BeNil())
		})
	})

	Describe("CreateToken", func() {
		It("创建新 Token", func() {
			token := CreateTestToken("new-user", models.PlainUser, false)
			err := tokenStore.CreateToken(token)
			Expect(err).ShouldNot(HaveOccurred())

			found := tokenStore.GetTokenByCondition(&models.BcsUser{Name: "new-user"})
			Expect(found).ShouldNot(BeNil())
			Expect(found.UserToken).Should(Equal("test-token-new-user"))
		})
	})

	Describe("UpdateToken", func() {
		It("更新 Token 信息", func() {
			token := CreateTestToken("update-user", models.PlainUser, false)
			err := tokenStore.CreateToken(token)
			Expect(err).ShouldNot(HaveOccurred())

			updatedToken := &models.BcsUser{
				UserToken: "updated-token",
				ExpiresAt: time.Now().Add(2 * time.Hour),
			}
			err = tokenStore.UpdateToken(token, updatedToken)
			Expect(err).ShouldNot(HaveOccurred())

			found := tokenStore.GetTokenByCondition(&models.BcsUser{Name: "update-user"})
			Expect(found).ShouldNot(BeNil())
		})
	})

	Describe("DeleteToken", func() {
		It("删除 Token（软删除）", func() {
			token := CreateTestToken("delete-user", models.PlainUser, false)
			err := tokenStore.CreateToken(token)
			Expect(err).ShouldNot(HaveOccurred())

			err = tokenStore.DeleteToken(token.UserToken)
			Expect(err).ShouldNot(HaveOccurred())

			found := tokenStore.GetTokenByCondition(&models.BcsUser{Name: "delete-user"})
			Expect(found).Should(BeNil())
		})
	})

	Describe("GetAllNotExpiredTokens", func() {
		It("获取所有未过期的 Token", func() {
			token := CreateTestToken("valid-user", models.PlainUser, false)
			err := tokenStore.CreateToken(token)
			Expect(err).ShouldNot(HaveOccurred())

			tokens := tokenStore.GetAllNotExpiredTokens()
			Expect(len(tokens)).Should(BeNumerically(">", 0))
		})
	})

	Describe("GetAllTokens", func() {
		It("获取所有 Token", func() {
			for i := 0; i < 3; i++ {
				token := CreateTestToken(fmt.Sprintf("all-user-%d", i), models.PlainUser, false)
				err := tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())
			}

			tokens := tokenStore.GetAllTokens()
			Expect(len(tokens)).Should(BeNumerically(">=", 3))
		})
	})
})
