# 集成测试文档

## 1. 概述

本文档描述 `bcs-user-manager` 项目的数据库集成测试方案，基于 Ginkgo v2 + Gomega 测试框架，支持 MySQL、达梦、OpenGauss 等数据库的自动化测试。

## 2. 测试环境准备

### 2.1 环境要求

- Go 1.24+
- MySQL 8.0+ (或使用 Docker)
- OpenSSL 1.1.1+（**重要**：SM4 加密需要 OpenSSL 1.1.1 以上版本，1.0.x 不支持）

### 2.2 CGO 环境变量

由于项目依赖 `TencentBlueKing/crypto-golang-sdk` (SM4 加密库)，**必须**设置 CGO 环境变量：

**macOS（Homebrew 安装 OpenSSL 3）：**

```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/openssl@3/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/openssl@3/lib -lssl -lcrypto"
```

**Linux（系统 OpenSSL 版本过低需自行编译安装）：**

如果系统 OpenSSL 版本低于 1.1.1（如 CentOS 7 自带的 1.0.2k），需要安装新版 OpenSSL：

```bash
# 下载并编译 OpenSSL 1.1.1（安装到 /opt/openssl，不影响系统）
cd /tmp
curl -LO https://www.openssl.org/source/openssl-1.1.1w.tar.gz
tar xzf openssl-1.1.1w.tar.gz
cd openssl-1.1.1w
./config --prefix=/opt/openssl
make -j$(nproc) && make install

# 设置环境变量（每次测试前执行，或写入 ~/.bashrc）
export CGO_CFLAGS="-I/opt/openssl/include"
export CGO_LDFLAGS="-L/opt/openssl/lib -lssl -lcrypto"
export LD_LIBRARY_PATH=/opt/openssl/lib:$LD_LIBRARY_PATH
```

建议将上述配置添加到 `~/.bashrc` 或 `~/.zshrc` 中。

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
    ├── common.go           # storeSet 定义、describeStoreTests() 入口
    ├── common_token.go     # Token Store 测试（16 方法）
    ├── common_user.go      # User Store 测试（3 方法）
    ├── common_credentials.go # Credentials Store 测试（7 方法）
    ├── common_register_token.go # RegisterToken Store 测试（2 方法）
    ├── common_cluster.go   # Cluster Store 测试（2 方法）
    ├── common_permission.go # Permission Store 测试（5 方法）
    ├── common_activity.go  # Activity Store 测试（3 方法）
    ├── common_log.go       # Log Store 测试（4 方法）
    ├── common_notify.go    # TokenNotify Store 测试（3 方法）
    ├── common_tke.go       # TkeCidr Store 测试（4 方法）
    ├── mysql_test.go      # MySQL 测试入口（默认）
    ├── dm_test.go         # 达梦测试入口
    └── gaussdb_test.go    # OpenGauss 测试入口
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
go test -tags=dm ./tests/sqlstore/... -v
```

### 4.3 OpenGauss 测试

```bash
# 需要先启动 OpenGauss 数据库
go test -tags=gaussdb ./tests/sqlstore/... -v
```

### 4.4 OceanBase IAM Migration 测试

`tests/iammigration` 用于验证 IAM 权限模型 migration 在 MySQL、OceanBase、GoldenDB、达梦和高斯数据库上的行为。各数据库入口仅负责连接初始化，migration 断言统一复用同一套测试逻辑。

```bash
BCS_TEST_MYSQL_HOST=<ob_host> \
BCS_TEST_MYSQL_PORT=<ob_port> \
BCS_TEST_MYSQL_USER=<ob_user> \
BCS_TEST_MYSQL_PASSWORD='<ob_password>' \
BCS_TEST_MYSQL_DBNAME=bcs_user_iam_migration_test \
BCS_TEST_DB_TYPE=ob \
go test ./tests/iammigration -run TestIAMMigrateMySQLCompatible -count=1 -v
```

MySQL 使用默认 `BCS_TEST_DB_TYPE=mysql`；GoldenDB 可使用 `dg` 或 `goldendb`；达梦和高斯分别使用：

```bash
go test -tags=dm ./tests/iammigration -run TestIAMMigrateDM -count=1 -v
go test -tags=gaussdb ./tests/iammigration -run TestIAMMigrateGaussDB -count=1 -v
```

测试会执行 `migrations/*.json` 中的 IAM migration，并验证版本表状态。它构造 `UserManager` 后调用服务启动同样使用的 `MigrateIAM(...)` 路径；测试使用本地 fake IAM server 响应 `/ping` 和模型接口，避免依赖真实 IAM 服务。

自动清理只会删除名称以默认测试前缀（`bcs_user_test`、`bcs_user_migration_test`、`bcs_user_iam_migration_test`）开头的数据库。自定义测试库名时，额外设置 `BCS_TEST_MYSQL_DBNAME_PREFIX=<custom_prefix>`；可用逗号分隔多个前缀。

## 5. OceanBase、GoldenDB 与 IAM Migration 适配说明

### 5.1 为什么需要适配

OceanBase 4.2 兼容 MySQL 协议，但不支持 MySQL advisory lock 函数：

```sql
SELECT GET_LOCK(?, 10)
SELECT RELEASE_LOCK(?)
```

`iam-go-sdk` 的默认 `Migrate(...)` 入口会在创建 `bk_iam_migrations` 版本表前执行上述锁操作。因此在 OceanBase 4.2 上，服务启动阶段的 IAM migration 会失败，错误类似：

```text
FUNCTION GET_LOCK does not exist
```

验证结果表明，跳过 advisory lock 后，`bk_iam_migrations` 表创建、版本写入，以及 `0000` 到 `0013` 的 migration 文件都可以在 OceanBase 上正常执行。因此 OceanBase 场景只需要跳过 lock，不需要修改 migration JSON 内容。

### 5.2 服务代码如何处理

业务表连接仍使用 MySQL driver。OceanBase 可将 `database_config.db_type` 配置为 `ob` 或 `oceanbase`，GoldenDB 可配置为 `dg` 或 `goldendb`。这两类数据库使用结构化 `database_config` 时，应省略 `mysql_dsn`（部署环境中的 `coreDatabaseDsn`）或将其保持为空。

- `sqlstore.InitCoreDatabase` 将 OceanBase 和 GoldenDB 按 MySQL 兼容数据库连接。
- IAM migration 使用 `iam-go-sdk` 的 `MigrateWithConfig(...)`。
- `iammigrate.Config.NoLock` 设置为 `true`，跳过 `GET_LOCK` / `RELEASE_LOCK`。

普通 MySQL 场景仍使用原来的 `u.IamPermClient.Migrate(...)`，保持原行为不变。

### 5.3 为什么修改 iam-go-sdk 依赖

当前 `go.mod` 保持：

```go
require github.com/TencentBlueKing/iam-go-sdk v0.1.6
```

同时通过 `replace` 指向内部适配版本：

```go
replace github.com/TencentBlueKing/iam-go-sdk => code.cwoa.net/rd-fy22-canway-platform-products/iam-go-sdk v0.1.5-alpha.0-xc.3
```

这样做的原因是：

- `bcs-common` 依赖解析仍保持 `github.com/TencentBlueKing/iam-go-sdk v0.1.6`，不会触发 `bcs-common` 降级。
- 内部 `v0.1.5-alpha.0-xc.3` 回补了 `MigrateWithConfig(...)` 和 `iammigrate.Config.NoLock`，可以解决 OceanBase 4.2 的 `GET_LOCK` 问题。
- 内部 `v0.1.5-alpha.0-xc.3` 的 `resource.Provider` 仍保持老接口，不需要为 `FetchInstanceList`、`FetchResourceTypeSchema` 增加额外兼容实现。

## 6. 测试覆盖

共覆盖 **10 个 Store**，**49 个方法**，**68 个测试用例**。

### 6.1 Token Store（16 方法）

| 测试用例 | 说明 |
|---------|------|
| `GetTokenByCondition - 按名称查询 Token` | 查询已存在的 Token |
| `GetTokenByCondition - 查询不存在的 Token 返回 nil` | 查询不存在时返回 nil |
| `CreateToken - 创建新 Token` | 创建 Token 并验证 |
| `UpdateToken - 更新 Token 信息` | 更新 Token 并验证 |
| `DeleteToken - 删除 Token（软删除）` | 软删除后查询返回 nil |
| `GetAllNotExpiredTokens - 获取所有未过期的 Token` | 获取未过期 Token |
| `GetAllTokens - 获取所有 Token` | 获取所有 Token |
| `GetUserTokensByName - 按用户名获取所有 Token（包括过期的）` | 同名用户多个 Token |
| `CreateTemporaryToken - 创建临时 Token` | 创建临时 Token |
| `GetTempTokenByCondition - 按 Token 查询临时 Token` | 按 Token 查询 |
| `GetTempTokenByCondition - 查询不存在的临时 Token 返回 nil` | 不存在返回 nil |
| `CreateClientToken - 创建客户端 Token` | 创建客户端 Token |
| `GetAllClients - 获取所有客户端` | 获取所有客户端 |
| `GetProjectClients - 获取指定项目的所有客户端` | 按项目过滤 |
| `GetClient - 获取指定项目和名称的客户端` | 获取单个客户端 |
| `UpdateClientToken - 更新客户端信息` | 更新客户端 |
| `DeleteProjectClient - 删除项目客户端` | 删除客户端 |

### 6.2 Credentials Store（7 方法）

| 测试用例 | 说明 |
|---------|------|
| `GetCredentials - 按 clusterId 查询凭证` | 按 clusterId 查询 |
| `GetCredentials - 查询不存在的凭证返回 nil` | 不存在返回 nil |
| `SaveCredentials - 创建凭证` | 创建新凭证 |
| `SaveCredentials - 更新已有凭证` | 覆盖更新 |
| `ListCredentials - 列出所有凭证` | 列表查询 |
| `SaveWsCredentials - 创建/更新 WebSocket 凭证` | 创建和更新 |
| `GetWsCredentials - 按 serverKey 查询` | 查询单个 |
| `GetWsCredentials - 查询不存在的 WebSocket 凭证返回 nil` | 不存在返回 nil |
| `DelWsCredentials - 删除 WebSocket 凭证` | 删除 |
| `GetWsCredentialsByClusterId - 按 clusterId 前缀查询` | 前缀匹配 |

### 6.3 RegisterToken Store（2 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateRegisterToken - 为集群创建注册 Token` | 创建并验证 |
| `CreateRegisterToken - 重复创建返回错误` | 唯一键冲突 |
| `GetRegisterToken - 按 clusterId 查询` | 查询单个 |
| `GetRegisterToken - 查询不存在的 clusterId 返回 nil` | 不存在返回 nil |

### 6.4 Cluster Store（2 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateCluster - 创建集群` | 创建并验证 |
| `GetCluster - 按 clusterId 查询集群` | 查询单个 |
| `GetCluster - 查询不存在的集群返回 nil` | 不存在返回 nil |

### 6.5 Permission Store（5 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateRole - 创建角色` | 创建并验证 |
| `GetRole - 按角色名查询` | 查询单个 |
| `GetRole - 查询不存在的角色返回 nil` | 不存在返回 nil |
| `CreateUserResourceRole - 创建用户资源角色关联` | 创建关联 |
| `GetUrrByCondition - 按条件查询关联` | 条件查询 |
| `GetUrrByCondition - 查询不存在的关联返回 nil` | 不存在返回 nil |
| `DeleteUserResourceRole - 删除关联` | 删除 |

### 6.6 Activity Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateActivity - 创建活动记录` | 创建单个 |
| `CreateActivity - 批量创建活动记录` | 批量创建 |
| `SearchActivities - 按项目编码搜索` | 按项目查 |
| `SearchActivities - 按资源类型过滤` | 资源类型过滤 |
| `SearchActivities - 按活动类型过滤` | 活动类型过滤 |
| `SearchActivities - projectCode 为空返回错误` | 参数校验 |
| `BatchDeleteActivity - 批量删除旧记录` | 按时间删除 |

### 6.7 Log Store（4 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateOperationLog - 创建操作日志` | 创建日志 |
| `ListOperationLogByClusterID - 按 clusterID 查询` | 按集群查 |
| `ListOperationLogByClusterID - 查询不存在的 clusterID 返回空` | 不存在返回空列表 |
| `ListOperationLogByUserClusterID - 按 clusterID 和用户名查询` | 双重条件 |
| `DeleteOperationLogByTime - 按时间范围删除` | 按时间删除 |

### 6.8 TokenNotify Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateTokenNotify - 创建通知记录` | 创建通知 |
| `GetTokenNotifyByCondition - 按 Token 查询` | 按 Token 查 |
| `GetTokenNotifyByCondition - 查询不存在的 Token 返回空列表` | 不存在返回空 |
| `DeleteTokenNotify - 删除通知记录` | 删除 |

### 6.9 TkeCidr Store（4 方法）

| 测试用例 | 说明 |
|---------|------|
| `SaveTkeCidr - 创建 TkeCidr` | 创建 CIDR |
| `QueryTkeCidr - 按 Vpc 和 Cidr 查询` | 联合查询 |
| `QueryTkeCidr - 查询不存在的记录返回 nil` | 不存在返回 nil |
| `UpdateTkeCidr - 更新 TkeCidr 信息` | 更新状态/集群 |
| `CountTkeCidr - 统计 TkeCidr` | 分组统计 |

### 6.10 User Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateUser - 创建用户` | 创建并验证 |
| `GetUserByCondition - 按名称查询用户` | 按名称查 |
| `GetUserByCondition - 查询不存在的用户返回 nil` | 不存在返回 nil |
| `UpdateUser - 更新用户信息` | 更新字段 |

## 7. 模型覆盖

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

## 8. 多数据库支持

### 8.1 Build Tags

| Tag | 说明 | 运行命令 |
|-----|------|---------|
| 无 tag（默认） | MySQL 测试 | `go test ./tests/sqlstore/... -v` |
| `dm` | 达梦测试 | `go test -tags=dm ./tests/sqlstore/... -v` |
| `gaussdb` | OpenGauss 测试 | `go test -tags=gaussdb ./tests/sqlstore/... -v` |

### 8.2 添加新测试

测试逻辑按 Store 拆分到 `common_*.go` 文件中：

1. 创建 `tests/sqlstore/common_xxx.go`，定义 `describeXxxTests(s *storeSet)` 函数
2. 在 `common.go` 的 `describeStoreTests()` 中调用 `describeXxxTests(s)`
3. 在对应 `common_xxx.go` 中添加测试用例

```go
// tests/sqlstore/common_user.go
func describeUserTests() {
    Describe("User Store 集成测试", func() {
        It("创建用户", func() { ... })
    })
}
```

```go
// tests/sqlstore/common.go
func describeStoreTests(s *storeSet) {
    describeTokenTests(s)
    describeUserTests()  // 新增
}
```

## 9. 注意事项

### 9.1 GORM AutoMigrate

使用 GORM AutoMigrate 自动创建表结构，无需手动编写 DDL，天然支持多数据库兼容。

### 9.2 全局 GCoreDB

部分 Store 函数依赖全局 `GCoreDB` 变量，测试框架通过 `sqlstore.SetGCoreDB(db)` 设置。

### 9.3 软删除

GORM 默认软删除，`DeleteToken` 等操作实际是更新 `deleted_at` 字段，查询时自动过滤。

### 9.4 OpenSSL 版本要求

**重要**：SM4 加密需要 OpenSSL 1.1.1+。项目依赖 `TencentBlueKing/crypto-golang-sdk`，其中 CGO 代码调用 `EVP_sm4_ctr`，该函数仅在 OpenSSL 1.1.1 及以上版本可用。

- **macOS**：Homebrew 安装的 `openssl@3` 已满足要求
- **Linux CentOS 7**：系统自带的 OpenSSL 1.0.2k 不支持，需要编译安装 OpenSSL 1.1.1+ 到 `/opt/openssl`（见 2.2 节）
- **其他 Linux 发行版**：如果 `openssl version` 输出低于 1.1.1，同样需要升级

### 9.5 调试技巧

- **详细测试步骤**：带上 `-args -ginkgo.v` 参数可以看到每个测试的详细步骤。
- **日志输出**：在测试中使用 `fmt.Fprintln(GinkgoWriter, "your message")` 输出调试信息，输出会与当前测试 Spec 绑定。

## 10. 测试输出示例

### 10.1 汇总模式 (-v)

```
=== RUN   TestTokenStoreMySQL
Running Suite: Token Store MySQL Integration Suite
======================================================================
Will run 68 of 68 Specs
Token Store MySQL 集成测试环境已就绪
••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••
Ran 68 of 68 Specs in 0.607 seconds
SUCCESS! -- 68 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestTokenStoreMySQL (0.61s)
PASS
```

### 10.2 详细模式 (-v -args -ginkgo.v)

```
Running Suite: Token Store MySQL Integration Suite
======================================================================
Running: go test -v ./tests/sqlstore/... -args -ginkgo.v
Rand Seed: 1776392475
Will run 68 of 68 specs

Token Store MySQL 集成测试环境已就绪

  [32m•[0m [32m•[0m [32m•[0m ...（68 个点）

  Token Store MySQL 集成测试
    GetTokenByCondition
      [32m✓[0m 按名称查询 Token (0.001s)
      [32m✓[0m 查询不存在的 Token 返回 nil (0.001s)
    CreateToken
      [32m✓[0m 创建新 Token (0.001s)
    ...
    TokenNotify Store 集成测试
      [32m✓[0m 创建 Token 通知记录 (0.001s)
      [32m✓[0m 按 Token 查询通知记录 (0.001s)
      [32m✓[0m 查询不存在的 Token 返回空列表 (0.001s)
      [32m✓[0m 删除 Token 通知记录 (0.001s)
    TkeCidr Store 集成测试
      [32m✓[0m 创建 TkeCidr (0.001s)
      ...
    User Store 集成测试
      [32m✓[0m 创建用户 (0.001s)
      [32m✓[0m 按名称查询用户 (0.001s)
      [32m✓[0m 查询不存在的用户返回 nil (0.001s)
      [32m✓[0m 更新用户信息 (0.001s)

  Ran 68 of 68 Specs in 0.607 seconds
  SUCCESS! -- 68 Passed | 0 Failed | 0 Pending | 0 Skipped
```
