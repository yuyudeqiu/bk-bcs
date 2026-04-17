package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describePermissionTests 定义了 Permission Store 的测试套件
func describePermissionTests() {
	Describe("Permission Store 集成测试", func() {
		Describe("Role", func() {
			Describe("CreateRole", func() {
				It("创建角色", func() {
					role := &models.BcsRole{
						Name:    "test-role-create",
						Actions: `["create","update"]`,
					}
					err := sqlstore.CreateRole(role)
					Expect(err).ShouldNot(HaveOccurred())

					found := sqlstore.GetRole(role.Name)
					Expect(found).ShouldNot(BeNil())
					Expect(found.Name).Should(Equal("test-role-create"))
				})
			})

			Describe("GetRole", func() {
				It("按角色名查询角色", func() {
					role := &models.BcsRole{
						Name:    "test-role-get",
						Actions: `["read"]`,
					}
					err := sqlstore.CreateRole(role)
					Expect(err).ShouldNot(HaveOccurred())

					found := sqlstore.GetRole(role.Name)
					Expect(found).ShouldNot(BeNil())
					Expect(found.Actions).Should(Equal(`["read"]`))
				})

				It("查询不存在的角色返回 nil", func() {
					found := sqlstore.GetRole("non-existent-role")
					Expect(found).Should(BeNil())
				})
			})
		})

		Describe("UserResourceRole", func() {
			Describe("CreateUserResourceRole", func() {
				It("创建用户资源角色关联", func() {
					role := &models.BcsRole{Name: "test-urr-role", Actions: `["admin"]`}
					err := sqlstore.CreateRole(role)
					Expect(err).ShouldNot(HaveOccurred())

					urr := &models.BcsUserResourceRole{
						UserId:       1,
						ResourceType: "cluster",
						Resource:     "cluster-1",
						RoleId:       role.ID,
					}
					err = sqlstore.CreateUserResourceRole(urr)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Describe("GetUrrByCondition", func() {
				It("按条件查询用户资源角色关联", func() {
					role := &models.BcsRole{Name: "test-urr-get-role", Actions: `["admin"]`}
					err := sqlstore.CreateRole(role)
					Expect(err).ShouldNot(HaveOccurred())

					urr := &models.BcsUserResourceRole{
						UserId:       2,
						ResourceType: "project",
						Resource:     "project-1",
						RoleId:       role.ID,
					}
					err = sqlstore.CreateUserResourceRole(urr)
					Expect(err).ShouldNot(HaveOccurred())

					found := sqlstore.GetUrrByCondition(&models.BcsUserResourceRole{UserId: 2, ResourceType: "project"})
					Expect(found).ShouldNot(BeNil())
					Expect(found.Resource).Should(Equal("project-1"))
				})

				It("查询不存在的关联返回 nil", func() {
					found := sqlstore.GetUrrByCondition(&models.BcsUserResourceRole{UserId: 999, ResourceType: "non-existent"})
					Expect(found).Should(BeNil())
				})
			})

			Describe("DeleteUserResourceRole", func() {
				It("删除用户资源角色关联", func() {
					role := &models.BcsRole{Name: "test-urr-del-role", Actions: `["admin"]`}
					err := sqlstore.CreateRole(role)
					Expect(err).ShouldNot(HaveOccurred())

					urr := &models.BcsUserResourceRole{
						UserId:       3,
						ResourceType: "cluster",
						Resource:     "cluster-del",
						RoleId:       role.ID,
					}
					err = sqlstore.CreateUserResourceRole(urr)
					Expect(err).ShouldNot(HaveOccurred())

					err = sqlstore.DeleteUserResourceRole(urr)
					Expect(err).ShouldNot(HaveOccurred())

					found := sqlstore.GetUrrByCondition(urr)
					Expect(found).Should(BeNil())
				})
			})
		})
	})
}
