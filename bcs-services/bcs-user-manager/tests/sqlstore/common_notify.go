package sqlstore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

// describeTokenNotifyTests 定义了 TokenNotify Store 的测试套件
func describeTokenNotifyTests(s *storeSet) {
	Describe("TokenNotify Store 集成测试", func() {
		Describe("CreateTokenNotify", func() {
			It("创建 Token 通知记录", func() {
				notify := &models.BcsTokenNotify{
					Token:      "test-token-notify-create",
					NotifyType: models.NotifyByEmail,
					Phase:      models.OverduePhase,
					Result:     true,
					Message:    "token expired",
					RequestID:  "req-001",
				}
				err := s.tokenNotifyStore.CreateTokenNotify(notify)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("GetTokenNotifyByCondition", func() {
			It("按 Token 查询通知记录", func() {
				notify := &models.BcsTokenNotify{
					Token:      "test-token-notify-get",
					NotifyType: models.NotifyByEmail,
					Phase:      models.DayPhase,
					Result:     false,
					RequestID:  "req-002",
				}
				err := s.tokenNotifyStore.CreateTokenNotify(notify)
				Expect(err).ShouldNot(HaveOccurred())

				results := s.tokenNotifyStore.GetTokenNotifyByCondition(&models.BcsTokenNotify{Token: "test-token-notify-get"})
				Expect(len(results)).Should(BeNumerically(">=", 1))
				Expect(results[0].NotifyType).Should(Equal(models.NotifyByEmail))
			})

			It("查询不存在的 Token 返回空列表", func() {
				results := s.tokenNotifyStore.GetTokenNotifyByCondition(&models.BcsTokenNotify{Token: "non-existent-token"})
				Expect(len(results)).Should(Equal(0))
			})
		})

		Describe("DeleteTokenNotify", func() {
			It("删除 Token 通知记录", func() {
				notify := &models.BcsTokenNotify{
					Token:      "test-token-notify-del",
					NotifyType: models.NotifyByRtx,
					Phase:      models.WeekPhase,
					Result:     true,
					RequestID:  "req-003",
				}
				err := s.tokenNotifyStore.CreateTokenNotify(notify)
				Expect(err).ShouldNot(HaveOccurred())

				err = s.tokenNotifyStore.DeleteTokenNotify("test-token-notify-del")
				Expect(err).ShouldNot(HaveOccurred())

				results := s.tokenNotifyStore.GetTokenNotifyByCondition(&models.BcsTokenNotify{Token: "test-token-notify-del"})
				Expect(len(results)).Should(Equal(0))
			})
		})
	})
}
