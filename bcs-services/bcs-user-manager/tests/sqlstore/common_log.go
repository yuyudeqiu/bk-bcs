package sqlstore_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeLogTests 定义了 Log Store 的测试套件
func describeLogTests() {
	Describe("Log Store 集成测试", func() {
		Describe("CreateOperationLog", func() {
			It("创建操作日志", func() {
				log := &models.BcsOperationLog{
					ClusterType: "testCluster",
					ClusterID:   "cluster-log-create",
					Path:        "/api/v1/cluster",
					Message:     "create cluster",
					OpUser:      "admin",
				}
				err := sqlstore.CreateOperationLog(log)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("ListOperationLogByClusterID", func() {
			It("按 clusterID 查询操作日志", func() {
				log := &models.BcsOperationLog{
					ClusterType: "testCluster",
					ClusterID:   "cluster-log-get",
					Path:        "/api/v1/cluster",
					Message:     "get cluster",
					OpUser:      "admin",
				}
				err := sqlstore.CreateOperationLog(log)
				Expect(err).ShouldNot(HaveOccurred())

				logs := sqlstore.ListOperationLogByClusterID("cluster-log-get")
				Expect(len(logs)).Should(BeNumerically(">=", 1))
			})

			It("查询不存在的 clusterID 返回空列表", func() {
				logs := sqlstore.ListOperationLogByClusterID("non-existent")
				Expect(len(logs)).Should(Equal(0))
			})
		})

		Describe("ListOperationLogByUserClusterID", func() {
			It("按 clusterID 和用户名查询操作日志", func() {
				log := &models.BcsOperationLog{
					ClusterType: "testCluster",
					ClusterID:   "cluster-log-user",
					Path:        "/api/v1/cluster",
					Message:     "update cluster",
					OpUser:      "testuser",
				}
				err := sqlstore.CreateOperationLog(log)
				Expect(err).ShouldNot(HaveOccurred())

				logs := sqlstore.ListOperationLogByUserClusterID("cluster-log-user", "testuser")
				Expect(len(logs)).Should(BeNumerically(">=", 1))
				Expect(logs[0].OpUser).Should(Equal("testuser"))
			})
		})

		Describe("DeleteOperationLogByTime", func() {
			It("按时间范围删除操作日志", func() {
				log := &models.BcsOperationLog{
					ClusterType: "testCluster",
					ClusterID:   "cluster-log-del",
					Path:        "/api/v1/cluster",
					Message:     "delete cluster",
					OpUser:      "admin",
				}
				err := sqlstore.CreateOperationLog(log)
				Expect(err).ShouldNot(HaveOccurred())

				// 删除 1 天前的记录
				err = sqlstore.DeleteOperationLogByTime(time.Now().Add(-24*time.Hour), time.Now().Add(-1*time.Hour))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
}
