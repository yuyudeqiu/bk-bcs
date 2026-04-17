package sqlstore_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

var zeroTime = time.Unix(1, 0)
var farFuture = time.Unix(2<<32, 0)

// describeActivityTests 定义了 Activity Store 的测试套件
func describeActivityTests() {
	Describe("Activity Store 集成测试", func() {
		Describe("CreateActivity", func() {
			It("创建活动记录", func() {
				activity := []*models.Activity{{
					ProjectCode:  "project-activity",
					ResourceType: "cluster",
					ResourceName: "test-cluster",
					ResourceID:   "cluster-1",
					ActivityType: "create",
					Status:       models.ActivityStatusSuccess,
					Username:     "admin",
					Description:  "创建集群",
				}}
				err := sqlstore.CreateActivity(activity)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("批量创建活动记录", func() {
				activities := []*models.Activity{
					{
						ProjectCode:  "project-batch",
						ResourceType: "cluster",
						ActivityType: "create",
						Status:       models.ActivityStatusSuccess,
						Username:     "admin",
					},
					{
						ProjectCode:  "project-batch",
						ResourceType: "cluster",
						ActivityType: "update",
						Status:       models.ActivityStatusSuccess,
						Username:     "admin",
					},
				}
				err := sqlstore.CreateActivity(activities)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("SearchActivities", func() {
			It("按项目编码搜索活动记录", func() {
				activity := []*models.Activity{{
					ProjectCode:  "project-search",
					ResourceType: "cluster",
					ActivityType: "create",
					Status:       models.ActivityStatusSuccess,
					Username:     "admin",
				}}
				err := sqlstore.CreateActivity(activity)
				Expect(err).ShouldNot(HaveOccurred())

				results, total, err := sqlstore.SearchActivities("project-search", "", "", 0, zeroTime, farFuture, 0, 10)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(total).Should(BeNumerically(">=", 1))
				Expect(len(results)).Should(BeNumerically(">=", 1))
			})

			It("按资源类型过滤搜索", func() {
				activity := []*models.Activity{{
					ProjectCode:  "project-filter",
					ResourceType: "cluster",
					ActivityType: "delete",
					Status:       models.ActivityStatusSuccess,
					Username:     "admin",
				}}
				err := sqlstore.CreateActivity(activity)
				Expect(err).ShouldNot(HaveOccurred())

				results, _, err := sqlstore.SearchActivities("project-filter", "cluster", "", 0, zeroTime, farFuture, 0, 10)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(results)).Should(BeNumerically(">=", 1))
			})

			It("按活动类型过滤搜索", func() {
				activity := []*models.Activity{{
					ProjectCode:  "project-at-filter",
					ResourceType: "cluster",
					ActivityType: "update",
					Status:       models.ActivityStatusSuccess,
					Username:     "admin",
				}}
				err := sqlstore.CreateActivity(activity)
				Expect(err).ShouldNot(HaveOccurred())

				results, _, err := sqlstore.SearchActivities("project-at-filter", "", "update", 0, zeroTime, farFuture, 0, 10)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(results)).Should(BeNumerically(">=", 1))
			})

			It("projectCode 为空返回错误", func() {
				_, _, err := sqlstore.SearchActivities("", "", "", 0, zeroTime, farFuture, 0, 10)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(ContainSubstring("projectCode can not be empty"))
			})
		})

		Describe("BatchDeleteActivity", func() {
			It("批量删除旧的活动记录", func() {
				activity := []*models.Activity{{
					ProjectCode:  "project-del",
					ResourceType: "cluster",
					ActivityType: "create",
					Status:       models.ActivityStatusSuccess,
					Username:     "admin",
				}}
				err := sqlstore.CreateActivity(activity)
				Expect(err).ShouldNot(HaveOccurred())

				// 删除 1 天前的记录
				err = sqlstore.BatchDeleteActivity([]string{"cluster"}, time.Now().Add(24*time.Hour))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
}
