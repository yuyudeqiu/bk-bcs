package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeClusterTests 定义了 Cluster Store 的测试套件
func describeClusterTests() {
	Describe("Cluster Store 集成测试", func() {
		Describe("CreateCluster", func() {
			It("创建集群", func() {
				cluster := &models.BcsCluster{
					ID:               "cluster-create",
					ClusterType:      1,
					TkeClusterId:     "tke-cluster-1",
					TkeClusterRegion: "ap-guangzhou",
				}
				err := sqlstore.CreateCluster(cluster)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.GetCluster(cluster.ID)
				Expect(found).ShouldNot(BeNil())
				Expect(found.ID).Should(Equal(cluster.ID))
			})
		})

		Describe("GetCluster", func() {
			It("按 clusterId 查询集群", func() {
				cluster := &models.BcsCluster{
					ID:               "cluster-get",
					ClusterType:      1,
					TkeClusterId:     "tke-cluster-get",
					TkeClusterRegion: "ap-guangzhou",
				}
				err := sqlstore.CreateCluster(cluster)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.GetCluster(cluster.ID)
				Expect(found).ShouldNot(BeNil())
				Expect(found.TkeClusterId).Should(Equal("tke-cluster-get"))
			})

			It("查询不存在的集群返回 nil", func() {
				found := sqlstore.GetCluster("non-existent-cluster")
				Expect(found).Should(BeNil())
			})
		})
	})
}
