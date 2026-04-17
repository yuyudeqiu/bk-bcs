package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeRegisterTokenTests 定义了 RegisterToken Store 的测试套件
func describeRegisterTokenTests() {
	Describe("RegisterToken Store 集成测试", func() {
		Describe("CreateRegisterToken", func() {
			It("为集群创建注册 Token", func() {
				clusterID := "cluster-register-create"
				err := sqlstore.CreateRegisterToken(clusterID)
				Expect(err).ShouldNot(HaveOccurred())

				token := sqlstore.GetRegisterToken(clusterID)
				Expect(token).ShouldNot(BeNil())
				Expect(token.ClusterId).Should(Equal(clusterID))
				Expect(token.Token).ShouldNot(BeEmpty())
			})

			It("重复创建同一集群的 Token 返回错误", func() {
				clusterID := "cluster-register-dup"
				err := sqlstore.CreateRegisterToken(clusterID)
				Expect(err).ShouldNot(HaveOccurred())

				err = sqlstore.CreateRegisterToken(clusterID)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(ContainSubstring("Duplicate entry"))
			})
		})

		Describe("GetRegisterToken", func() {
			It("按 clusterId 查询注册 Token", func() {
				clusterID := "cluster-register-get"
				err := sqlstore.CreateRegisterToken(clusterID)
				Expect(err).ShouldNot(HaveOccurred())

				token := sqlstore.GetRegisterToken(clusterID)
				Expect(token).ShouldNot(BeNil())
				Expect(token.ClusterId).Should(Equal(clusterID))
			})

			It("查询不存在的 clusterId 返回 nil", func() {
				token := sqlstore.GetRegisterToken("non-existent-cluster")
				Expect(token).Should(BeNil())
			})
		})
	})
}
