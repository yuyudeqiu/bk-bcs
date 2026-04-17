package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// describeCredentialsTests 定义了 Credentials Store 的测试套件
func describeCredentialsTests() {
	Describe("Credentials Store 集成测试", func() {
		Describe("GetCredentials", func() {
			It("按 clusterId 查询凭证", func() {
				clusterID := "cluster-credentials-get"
				err := sqlstore.SaveCredentials(clusterID, "127.0.0.1:6443", "ca-data", "user-token", "cluster.domain")
				Expect(err).ShouldNot(HaveOccurred())

				cred := sqlstore.GetCredentials(clusterID)
				Expect(cred).ShouldNot(BeNil())
				Expect(cred.ClusterId).Should(Equal(clusterID))
				Expect(cred.ServerAddresses).Should(Equal("127.0.0.1:6443"))
			})

			It("查询不存在的凭证返回 nil", func() {
				cred := sqlstore.GetCredentials("non-existent-cluster")
				Expect(cred).Should(BeNil())
			})
		})

		Describe("SaveCredentials", func() {
			It("创建凭证", func() {
				clusterID := "cluster-credentials-create"
				err := sqlstore.SaveCredentials(clusterID, "127.0.0.1:6443", "ca-data", "user-token", "cluster.domain")
				Expect(err).ShouldNot(HaveOccurred())

				cred := sqlstore.GetCredentials(clusterID)
				Expect(cred).ShouldNot(BeNil())
			})

			It("更新已有凭证", func() {
				clusterID := "cluster-credentials-update"
				err := sqlstore.SaveCredentials(clusterID, "127.0.0.1:6443", "ca-data", "user-token", "cluster.domain")
				Expect(err).ShouldNot(HaveOccurred())

				err = sqlstore.SaveCredentials(clusterID, "127.0.0.1:7443", "new-ca-data", "new-token", "new.cluster.domain")
				Expect(err).ShouldNot(HaveOccurred())

				cred := sqlstore.GetCredentials(clusterID)
				Expect(cred).ShouldNot(BeNil())
				Expect(cred.ServerAddresses).Should(Equal("127.0.0.1:7443"))
				Expect(cred.UserToken).Should(Equal("new-token"))
			})
		})

		Describe("ListCredentials", func() {
			It("列出所有凭证", func() {
				clusterID := "cluster-credentials-list"
				err := sqlstore.SaveCredentials(clusterID, "127.0.0.1:6443", "ca-data", "user-token", "cluster.domain")
				Expect(err).ShouldNot(HaveOccurred())

				creds := sqlstore.ListCredentials()
				Expect(len(creds)).Should(BeNumerically(">", 0))
			})
		})

		Describe("WebSocket Credentials", func() {
			Describe("SaveWsCredentials", func() {
				It("创建 WebSocket 凭证", func() {
					serverKey := "ws-key-save"
					err := sqlstore.SaveWsCredentials(serverKey, "client-module", "server-address", "ca-data", "user-token")
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("更新已有 WebSocket 凭证", func() {
					serverKey := "ws-key-update"
					err := sqlstore.SaveWsCredentials(serverKey, "client-module", "server-address", "ca-data", "user-token")
					Expect(err).ShouldNot(HaveOccurred())

					err = sqlstore.SaveWsCredentials(serverKey, "new-module", "new-address", "new-ca", "new-token")
					Expect(err).ShouldNot(HaveOccurred())

					cred := sqlstore.GetWsCredentials(serverKey)
					Expect(cred).ShouldNot(BeNil())
					Expect(cred.ClientModule).Should(Equal("new-module"))
				})
			})

			Describe("GetWsCredentials", func() {
				It("按 serverKey 查询 WebSocket 凭证", func() {
					serverKey := "ws-key-get"
					err := sqlstore.SaveWsCredentials(serverKey, "client-module", "server-address", "ca-data", "user-token")
					Expect(err).ShouldNot(HaveOccurred())

					cred := sqlstore.GetWsCredentials(serverKey)
					Expect(cred).ShouldNot(BeNil())
					Expect(cred.ServerKey).Should(Equal(serverKey))
				})

				It("查询不存在的 WebSocket 凭证返回 nil", func() {
					cred := sqlstore.GetWsCredentials("non-existent-key")
					Expect(cred).Should(BeNil())
				})
			})

			Describe("DelWsCredentials", func() {
				It("删除 WebSocket 凭证", func() {
					serverKey := "ws-key-del"
					err := sqlstore.SaveWsCredentials(serverKey, "client-module", "server-address", "ca-data", "user-token")
					Expect(err).ShouldNot(HaveOccurred())

					sqlstore.DelWsCredentials(serverKey)

					cred := sqlstore.GetWsCredentials(serverKey)
					Expect(cred).Should(BeNil())
				})
			})

			Describe("GetWsCredentialsByClusterId", func() {
				It("按 clusterId 前缀查询 WebSocket 凭证", func() {
					clusterID := "ws-cluster-query"
					serverKey := clusterID + "-node-1"
					err := sqlstore.SaveWsCredentials(serverKey, "client-module", "server-address", "ca-data", "user-token")
					Expect(err).ShouldNot(HaveOccurred())

					creds := sqlstore.GetWsCredentialsByClusterId(clusterID)
					Expect(len(creds)).Should(BeNumerically(">=", 1))
				})
			})
		})
	})
}
