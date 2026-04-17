package sqlstore_test

import (
	"gorm.io/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
)

// storeSet 包含所有 store 实例
type storeSet struct {
	db         *gorm.DB
	tokenStore sqlstore.TokenStore
	// 后续新增其他 store...
}

// describeStoreTests 定义了通用的测试套件入口
// 在各数据库的 _test.go 中调用
func describeStoreTests(s *storeSet) {
	// 测试 Token 相关的 CRUD
	describeTokenTests(s)
	// 后续新增其他 store 的测试...
}
