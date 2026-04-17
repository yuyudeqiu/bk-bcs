# 集成测试文档

## 1. 概述

本文档描述 `bcs-user-manager` 项目的数据库集成测试方案，基于 Ginkgo v2 + Gomega 测试框架，支持 MySQL、达梦、OpenGauss 等数据库的自动化测试。

## 2. 测试环境准备

### 2.1 环境要求

- Go 1.24+
- MySQL 8.0+ (或使用 Docker)
- OpenSSL

### 2.2 环境变量

由于项目依赖 `TencentBlueKing/crypto-golang-sdk` (SM4 加密库)，需要设置 CGO 环境变量：

```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/openssl@3/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/openssl@3/lib -lssl -lcrypto"
```

建议将上述配置添加到 `~/.zshrc` 中。

### 2.3 运行测试

```bash
go test ./tests/sqlstore/... -v
```

## 3. 测试结构

### 3.1 目录结构

```
tests/
├── framework/
│   ├── common.go       # 公共函数：NewDBClient, InitTables, CleanData
│   ├── mysql.go        # MySQL 初始化
│   ├── dm.go           # 达梦初始化
│   └── gaussdb.go      # OpenGauss 初始化
└── sqlstore/
    ├── mysql_test.go    # MySQL 测试 (默认)
    ├── dm_test.go       # 达梦测试
    └── gaussdb_test.go  # OpenGauss 测试
```

### 3.2 测试框架 (tests/framework/)

| 文件 | 函数 | 说明 |
|------|------|------|
| `common.go` | `NewDBClient(cfg)` | 使用 SDK 创建数据库客户端 |
| `common.go` | `InitTables(db)` | 使用 GORM AutoMigrate 创建表结构 |
| `common.go` | `CleanData(db)` | 清理测试数据 |
| `mysql.go` | `InitMySQL(dbName)` | 初始化 MySQL 测试数据库 |
| `dm.go` | `InitDM(dbName)` | 初始化达梦测试数据库 |
| `gaussdb.go` | `InitGaussDB(dbName)` | 初始化 OpenGauss 测试数据库 |

### 3.3 数据库配置

各数据库的默认配置在 `framework/*Config` 变量中：

| 数据库 | 配置文件 | 默认 Host | 默认 Port | 默认用户 | 默认密码 |
|--------|---------|-----------|-----------|----------|----------|
| MySQL | `MySQLConfig` | localhost | 3306 | root | root |
| 达梦 | `DMConfig` | localhost | 5236 | bkauth | bkauth123 |
| OpenGauss | `GaussDBConfig` | localhost | 5433 | gaussdb | openGauss@123 |

## 4. 运行测试

### 4.1 MySQL 测试（默认）

```bash
# 运行测试（显示汇总）
go test ./tests/sqlstore/... -v

# 显示详细测试步骤
go test ./tests/sqlstore/... -v -args -ginkgo.v
```

### 4.2 达梦测试

```bash
# 需要先启动达梦数据库
go test -tags=dm ./tests/sqlstore/... -v -args -ginkgo.v
```

### 4.3 OpenGauss 测试

```bash
# 需要先启动 OpenGauss 数据库
go test -tags=gaussdb ./tests/sqlstore/... -v -args -ginkgo.v
```

## 5. 测试覆盖

### 5.1 Token Store 测试用例

| 测试用例 | 说明 |
|---------|------|
| `GetTokenByCondition - 按名称查询 Token` | 查询已存在的 Token |
| `GetTokenByCondition - 查询不存在的 Token 返回 nil` | 查询不存在时返回 nil |
| `CreateToken - 创建新 Token` | 创建 Token 并验证 |
| `UpdateToken - 更新 Token 信息` | 更新 Token 并验证 |
| `DeleteToken - 删除 Token（软删除）` | 软删除后查询返回 nil |
| `GetAllNotExpiredTokens - 获取所有未过期的 Token` | 获取未过期 Token |
| `GetAllTokens - 获取所有 Token` | 获取所有 Token |

## 6. 模型覆盖

测试通过 GORM AutoMigrate 覆盖以下模型：

| 模型 | 所在文件 |
|------|---------|
| `BcsUser` | `models/user.go` |
| `BcsTempToken` | `models/token.go` |
| `BcsClient` | `models/user.go` |
| `BcsClientUser` | `models/user.go` |
| `BcsCluster` | `models/cluster.go` |
| `BcsRegisterToken` | `models/cluster.go` |
| `BcsClusterCredential` | `models/cluster.go` |
| `BcsWsClusterCredentials` | `models/cluster.go` |
| `BcsRole` | `models/permission.go` |
| `BcsUserResourceRole` | `models/permission.go` |
| `Activity` | `models/activity.go` |
| `BcsOperationLog` | `models/log.go` |
| `BcsTokenNotify` | `models/notify.go` |
| `TkeCidr` | `models/tke.go` |

## 7. 多数据库支持

### 7.1 Build Tags

| Tag | 说明 | 运行命令 |
|-----|------|---------|
| 无 tag（默认） | MySQL 测试 | `go test ./tests/...` |
| `dm` | 达梦测试 | `go test -tags=dm ./tests/...` |
| `gaussdb` | OpenGauss 测试 | `go test -tags=gaussdb ./tests/...` |

### 7.2 扩展测试到其他 Store

参考 `tests/sqlstore/mysql_test.go` 的模式，创建新的测试文件：

```go
// mysql_test.go
//go:build !dm && !gaussdb
package sqlstore_test

var _ = Describe("Token Store MySQL 集成测试", func() {
    // 测试用例
})
```

## 8. 注意事项

### 8.1 GORM AutoMigrate

使用 GORM AutoMigrate 自动创建表结构，无需手动编写 DDL，天然支持多数据库兼容。

### 8.2 全局 GCoreDB

部分 Store 函数依赖全局 `GCoreDB` 变量，测试框架通过 `sqlstore.SetGCoreDB(db)` 设置。

### 8.3 软删除

GORM 默认软删除，`DeleteToken` 等操作实际是更新 `deleted_at` 字段，查询时自动过滤。

### 8.4 OpenSSL 依赖

由于项目依赖 `TencentBlueKing/crypto-golang-sdk` (SM4 加密库)，需要设置 CGO 环境变量指定 OpenSSL 路径：

```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/openssl@3/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/openssl@3/lib -lssl -lcrypto"
```

### 8.5 调试技巧

- **详细测试步骤**：带上 `-args -ginkgo.v` 参数可以看到每个测试的详细步骤。
- **日志输出**：在测试中使用 `fmt.Fprintln(GinkgoWriter, "your message")` 输出调试信息，输出会与当前测试 Spec 绑定。

## 9. 测试输出示例

### 9.1 汇总模式 (-v)

```
=== RUN   TestTokenStoreMySQL
Running Suite: Token Store MySQL Integration Suite
======================================================================
Will run 7 of 7 Specs
Token Store MySQL 集成测试环境已就绪
•••••••
Ran 7 of 7 Specs in 0.164 seconds
SUCCESS! -- 7 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestTokenStoreMySQL (0.16s)
PASS
```

### 9.2 详细模式 (-v -args -ginkgo.v)

```
Running Suite: Token Store MySQL Integration Suite
======================================================================
Running: go test -v ./tests/sqlstore/... -args -ginkgo.v
Rand Seed: 1776392475
Will run 7 of 7 specs

Token Store MySQL 集成测试环境已就绪

  [32m•[0m [32m•[0m [32m•[0m [32m•[0m [32m•[0m [32m•[0m [32m•[0m

  Token Store MySQL 集成测试
    GetTokenByCondition
      [32m✓[0m 按名称查询 Token (0.001s)
      [32m✓[0m 查询不存在的 Token 返回 nil (0.001s)
    CreateToken
      [32m✓[0m 创建新 Token (0.001s)
    UpdateToken
      [32m✓[0m 更新 Token 信息 (0.001s)
    DeleteToken
      [32m✓[0m 删除 Token（软删除） (0.001s)
    GetAllNotExpiredTokens
      [32m✓[0m 获取所有未过期的 Token (0.001s)
    GetAllTokens
      [32m✓[0m 获取所有 Token (0.001s)

  Ran 7 of 7 Specs in 0.164 seconds
  SUCCESS! -- 7 Passed | 0 Failed | 0 Pending | 0 Skipped
```
