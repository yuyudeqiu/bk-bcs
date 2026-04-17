package sqlstore_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
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

// CreateTestTempToken 创建测试临时 Token
func CreateTestTempToken(username string, expired bool) *models.BcsTempToken {
	expiresAt := time.Now().Add(time.Hour)
	if expired {
		expiresAt = time.Now().Add(-time.Hour)
	}
	return &models.BcsTempToken{
		Username:  username,
		Token:     "test-temp-token-" + username,
		UserType:  models.PlainUser,
		CreatedBy: "system",
		ExpiresAt: expiresAt,
	}
}

// ptrString returns pointer to string
func ptrString(s string) *string {
	return &s
}

// CreateTestClientUser 创建测试客户端用户
func CreateTestClientUser(projectCode, name string, expired bool) *models.BcsClientUser {
	expiresAt := time.Now().Add(time.Hour)
	if expired {
		expiresAt = time.Now().Add(-time.Hour)
	}
	return &models.BcsClientUser{
		ProjectCode:   projectCode,
		Name:          name,
		UserType:      models.PlainUser,
		UserToken:     "test-client-token-" + name,
		CreatedBy:     "system",
		Manager:       "admin",
		AuthorityUser: "admin",
		ExpiresAt:     expiresAt,
	}
}

// describeTokenTests 定义了 Token Store 的测试套件
func describeTokenTests(s *storeSet) {
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

		Describe("GetUserTokensByName", func() {
			It("按用户名获取所有 Token（包括过期的）", func() {
				name := "multi-token-user"
				token1 := CreateTestToken(name, models.PlainUser, false)
				token1.UserToken = "test-token-multi-1"
				token2 := CreateTestToken(name, models.PlainUser, true)
				token2.UserToken = "test-token-multi-2"
				err := s.tokenStore.CreateToken(token1)
				Expect(err).ShouldNot(HaveOccurred())
				err = s.tokenStore.CreateToken(token2)
				Expect(err).ShouldNot(HaveOccurred())

				tokens := s.tokenStore.GetUserTokensByName(name)
				Expect(len(tokens)).Should(Equal(2))
			})
		})

		Describe("Temporary Token", func() {
			Describe("CreateTemporaryToken", func() {
				It("创建临时 Token", func() {
					tempToken := CreateTestTempToken("temp-user", false)
					err := s.tokenStore.CreateTemporaryToken(tempToken)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Describe("GetTempTokenByCondition", func() {
				It("按 Token 查询临时 Token", func() {
					tempToken := CreateTestTempToken("temp-query-user", false)
					err := s.tokenStore.CreateTemporaryToken(tempToken)
					Expect(err).ShouldNot(HaveOccurred())

					found := s.tokenStore.GetTempTokenByCondition(&models.BcsTempToken{Token: tempToken.Token})
					Expect(found).ShouldNot(BeNil())
					Expect(found.Username).Should(Equal("temp-query-user"))
				})

				It("查询不存在的临时 Token 返回 nil", func() {
					found := s.tokenStore.GetTempTokenByCondition(&models.BcsTempToken{Token: "non-existent"})
					Expect(found).Should(BeNil())
				})
			})
		})

		Describe("Client Token", func() {
			Describe("CreateClientToken", func() {
				It("创建客户端 Token", func() {
					clientUser := CreateTestClientUser("project-a", "client-user-1", false)
					err := s.tokenStore.CreateClientToken(clientUser)
					Expect(err).ShouldNot(HaveOccurred())

					client := s.tokenStore.GetClient("project-a", "client-user-1")
					Expect(client).ShouldNot(BeNil())
					Expect(client.ProjectCode).Should(Equal("project-a"))
				})
			})

			Describe("GetAllClients", func() {
				It("获取所有客户端", func() {
					client := CreateTestClientUser("project-b", "client-all-1", false)
					err := s.tokenStore.CreateClientToken(client)
					Expect(err).ShouldNot(HaveOccurred())

					clients := s.tokenStore.GetAllClients()
					Expect(len(clients)).Should(BeNumerically(">", 0))
				})
			})

			Describe("GetProjectClients", func() {
				It("获取指定项目的所有客户端", func() {
					client := CreateTestClientUser("project-c", "client-proj-1", false)
					err := s.tokenStore.CreateClientToken(client)
					Expect(err).ShouldNot(HaveOccurred())

					clients := s.tokenStore.GetProjectClients("project-c")
					Expect(len(clients)).Should(BeNumerically(">=", 1))
				})
			})

			Describe("GetClient", func() {
				It("获取指定项目和名称的客户端", func() {
					client := CreateTestClientUser("project-d", "client-get-1", false)
					err := s.tokenStore.CreateClientToken(client)
					Expect(err).ShouldNot(HaveOccurred())

					found := s.tokenStore.GetClient("project-d", "client-get-1")
					Expect(found).ShouldNot(BeNil())
					Expect(found.Name).Should(Equal("client-get-1"))
				})

				It("查询不存在的客户端返回空对象", func() {
					found := s.tokenStore.GetClient("non-existent", "non-existent")
					Expect(found).ShouldNot(BeNil())
					// GetClient uses Raw().Scan() which returns empty struct, not nil
				})
			})

			Describe("UpdateClientToken", func() {
				It("更新客户端信息", func() {
					client := CreateTestClientUser("project-e", "client-update-1", false)
					err := s.tokenStore.CreateClientToken(client)
					Expect(err).ShouldNot(HaveOccurred())

					updatedClient := &models.BcsClient{
						Manager: ptrString("new-admin"),
					}
					err = s.tokenStore.UpdateClientToken("project-e", "client-update-1", updatedClient)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Describe("DeleteProjectClient", func() {
				It("删除项目客户端", func() {
					client := CreateTestClientUser("project-f", "client-del-1", false)
					err := s.tokenStore.CreateClientToken(client)
					Expect(err).ShouldNot(HaveOccurred())

					err = s.tokenStore.DeleteProjectClient("project-f", "client-del-1")
					Expect(err).ShouldNot(HaveOccurred())

					found := s.tokenStore.GetClient("project-f", "client-del-1")
					// 软删除后查不到，但返回空 struct
					Expect(found).ShouldNot(BeNil())
				})
			})
		})
	})
}
