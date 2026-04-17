package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeTkeCidrTests 定义了 TkeCidr Store 的测试套件
func describeTkeCidrTests() {
	Describe("TkeCidr Store 集成测试", func() {
		Describe("SaveTkeCidr", func() {
			It("创建 TkeCidr", func() {
				vpc := "vpc-tke-save"
				cidr := "10.0.0.0/16"
				cluster := "cluster-tke-1"
				err := sqlstore.SaveTkeCidr(vpc, cidr, 65536, sqlstore.CidrStatusAvailable, cluster)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.QueryTkeCidr(&models.TkeCidr{Vpc: vpc, Cidr: cidr})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Vpc).Should(Equal(vpc))
			})
		})

		Describe("QueryTkeCidr", func() {
			It("按 Vpc 和 Cidr 查询", func() {
				vpc := "vpc-tke-query"
				cidr := "10.1.0.0/16"
				err := sqlstore.SaveTkeCidr(vpc, cidr, 32768, sqlstore.CidrStatusUsed, "cluster-tke-query")
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.QueryTkeCidr(&models.TkeCidr{Vpc: vpc, Cidr: cidr})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Status).Should(Equal(sqlstore.CidrStatusUsed))
			})

			It("查询不存在的记录返回 nil", func() {
				found := sqlstore.QueryTkeCidr(&models.TkeCidr{Vpc: "non-existent-vpc", Cidr: "0.0.0.0/0"})
				Expect(found).Should(BeNil())
			})
		})

		Describe("UpdateTkeCidr", func() {
			It("更新 TkeCidr 信息", func() {
				vpc := "vpc-tke-update"
				cidr := "10.2.0.0/16"
				err := sqlstore.SaveTkeCidr(vpc, cidr, 16384, sqlstore.CidrStatusAvailable, "cluster-old")
				Expect(err).ShouldNot(HaveOccurred())

				original := sqlstore.QueryTkeCidr(&models.TkeCidr{Vpc: vpc, Cidr: cidr})
				Expect(original).ShouldNot(BeNil())

				updated := &models.TkeCidr{
					Status:  sqlstore.CidrStatusUsed,
					Cluster: ptrString("cluster-new"),
				}
				err = sqlstore.UpdateTkeCidr(original, updated)
				Expect(err).ShouldNot(HaveOccurred())

				found := sqlstore.QueryTkeCidr(&models.TkeCidr{Vpc: vpc, Cidr: cidr})
				Expect(found).ShouldNot(BeNil())
				Expect(found.Status).Should(Equal(sqlstore.CidrStatusUsed))
			})
		})

		Describe("CountTkeCidr", func() {
			It("统计 TkeCidr", func() {
				vpc := "vpc-tke-count"
				err := sqlstore.SaveTkeCidr(vpc, "10.3.0.0/16", 8192, sqlstore.CidrStatusReserved, "cluster-count")
				Expect(err).ShouldNot(HaveOccurred())

				counts := sqlstore.CountTkeCidr()
				Expect(len(counts)).Should(BeNumerically(">=", 1))
			})
		})
	})
}
