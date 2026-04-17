package sqlstore_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeUserTests 定义了 User Store 的测试套件
func describeUserTests() {
	Describe("User Store 集成测试", func() {
		Describe("CreateUser", func() {
			It("创建用户", func() {
				user := &models.BcsUser{
					Name:      "test-user-create",
					UserType:  models.PlainUser,
					UserToken: "test-user-token-create-" + uniqueID(),
					CreatedBy: "system",
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				err := sqlstore.CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.GetUserByCondition(&models.BcsUser{Name: "test-user-create"})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Name).Should(Equal("test-user-create"))
				Expect(found.UserType).Should(BeEquivalentTo(models.PlainUser))
			})
		})

		Describe("GetUserByCondition", func() {
			It("按名称查询用户", func() {
				user := &models.BcsUser{
					Name:      "test-user-get",
					UserType:  models.AdminUser,
					UserToken: "test-user-token-get-" + uniqueID(),
					CreatedBy: "system",
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				err := sqlstore.CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.GetUserByCondition(&models.BcsUser{Name: "test-user-get"})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Name).Should(Equal("test-user-get"))
			})

			It("查询不存在的用户返回 nil", func() {
				found := sqlstore.GetUserByCondition(&models.BcsUser{Name: "non-existent-user"})
				Expect(found).Should(BeNil())
			})
		})

		Describe("UpdateUser", func() {
			It("更新用户信息", func() {
				user := &models.BcsUser{
					Name:      "test-user-update",
					UserType:  models.PlainUser,
					UserToken: "test-user-token-update-" + uniqueID(),
					CreatedBy: "system",
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				err := sqlstore.CreateUser(user)
				Expect(err).ShouldNot(HaveOccurred())

				updated := &models.BcsUser{
					ExpiresAt: time.Now().Add(48 * time.Hour),
				}
				err = sqlstore.UpdateUser(user, updated)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.GetUserByCondition(&models.BcsUser{Name: "test-user-update"})
				Expect(found).ShouldNot(BeNil())
				Expect(found.ExpiresAt.After(time.Now().Add(24 * time.Hour))).Should(BeTrue())
			})
		})
	})
}

// uniqueID 返回唯一 ID，用于避免 UserToken unique 冲突
func uniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
