package sqlstore_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// storeSet 包含所有 store 实例
type storeSet struct {
	db         *gorm.DB
	tokenStore sqlstore.TokenStore
}

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

// describeStoreTests 定义了通用的测试套件
func describeStoreTests(s *storeSet) {
	Describe("Token Store 集成测试", func() {
		Describe("GetTokenByCondition", func() {
			It("按名称查询 Token", func() {
				token := CreateTestToken("test-user", models.PlainUser, false)
				err := s.tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())

				found := s.tokenStore.GetTokenByCondition(&models.BcsUser{Name: "test-user"})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Name).Should(Equal("test-user"))
			})

			It("查询不存在的 Token 返回 nil", func() {
				found := s.tokenStore.GetTokenByCondition(&models.BcsUser{Name: "non-existent"})
				Expect(found).Should(BeNil())
			})
		})

		Describe("CreateToken", func() {
			It("创建新 Token", func() {
				token := CreateTestToken("new-user", models.PlainUser, false)
				err := s.tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())

				found := s.tokenStore.GetTokenByCondition(&models.BcsUser{Name: "new-user"})
				Expect(found).ShouldNot(BeNil())
				Expect(found.UserToken).Should(Equal("test-token-new-user"))
			})
		})

		Describe("UpdateToken", func() {
			It("更新 Token 信息", func() {
				token := CreateTestToken("update-user", models.PlainUser, false)
				err := s.tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())

				updatedToken := &models.BcsUser{
					UserToken: "updated-token",
					ExpiresAt: time.Now().Add(2 * time.Hour),
				}
				err = s.tokenStore.UpdateToken(token, updatedToken)
				Expect(err).ShouldNot(HaveOccurred())

				found := s.tokenStore.GetTokenByCondition(&models.BcsUser{Name: "update-user"})
				Expect(found).ShouldNot(BeNil())
			})
		})

		Describe("DeleteToken", func() {
			It("删除 Token（软删除）", func() {
				token := CreateTestToken("delete-user", models.PlainUser, false)
				err := s.tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())

				err = s.tokenStore.DeleteToken(token.UserToken)
				Expect(err).ShouldNot(HaveOccurred())

				found := s.tokenStore.GetTokenByCondition(&models.BcsUser{Name: "delete-user"})
				Expect(found).Should(BeNil())
			})
		})

		Describe("GetAllNotExpiredTokens", func() {
			It("获取所有未过期的 Token", func() {
				token := CreateTestToken("valid-user", models.PlainUser, false)
				err := s.tokenStore.CreateToken(token)
				Expect(err).ShouldNot(HaveOccurred())

				tokens := s.tokenStore.GetAllNotExpiredTokens()
				Expect(len(tokens)).Should(BeNumerically(">", 0))
			})
		})

		Describe("GetAllTokens", func() {
			It("获取所有 Token", func() {
				for i := 0; i < 3; i++ {
					token := CreateTestToken(fmt.Sprintf("all-user-%d", i), models.PlainUser, false)
					err := s.tokenStore.CreateToken(token)
					Expect(err).ShouldNot(HaveOccurred())
				}

				tokens := s.tokenStore.GetAllTokens()
				Expect(len(tokens)).Should(BeNumerically(">=", 3))
			})
		})
	})
}
