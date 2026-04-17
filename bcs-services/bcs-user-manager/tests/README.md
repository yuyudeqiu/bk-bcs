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

## 5. 测试覆盖

共覆盖 **10 个 Store**，**49 个方法**，**68 个测试用例**。

### 5.1 Token Store（16 方法）

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

### 5.2 Credentials Store（7 方法）

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

### 5.3 RegisterToken Store（2 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateRegisterToken - 为集群创建注册 Token` | 创建并验证 |
| `CreateRegisterToken - 重复创建返回错误` | 唯一键冲突 |
| `GetRegisterToken - 按 clusterId 查询` | 查询单个 |
| `GetRegisterToken - 查询不存在的 clusterId 返回 nil` | 不存在返回 nil |

### 5.4 Cluster Store（2 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateCluster - 创建集群` | 创建并验证 |
| `GetCluster - 按 clusterId 查询集群` | 查询单个 |
| `GetCluster - 查询不存在的集群返回 nil` | 不存在返回 nil |

### 5.5 Permission Store（5 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateRole - 创建角色` | 创建并验证 |
| `GetRole - 按角色名查询` | 查询单个 |
| `GetRole - 查询不存在的角色返回 nil` | 不存在返回 nil |
| `CreateUserResourceRole - 创建用户资源角色关联` | 创建关联 |
| `GetUrrByCondition - 按条件查询关联` | 条件查询 |
| `GetUrrByCondition - 查询不存在的关联返回 nil` | 不存在返回 nil |
| `DeleteUserResourceRole - 删除关联` | 删除 |

### 5.6 Activity Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateActivity - 创建活动记录` | 创建单个 |
| `CreateActivity - 批量创建活动记录` | 批量创建 |
| `SearchActivities - 按项目编码搜索` | 按项目查 |
| `SearchActivities - 按资源类型过滤` | 资源类型过滤 |
| `SearchActivities - 按活动类型过滤` | 活动类型过滤 |
| `SearchActivities - projectCode 为空返回错误` | 参数校验 |
| `BatchDeleteActivity - 批量删除旧记录` | 按时间删除 |

### 5.7 Log Store（4 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateOperationLog - 创建操作日志` | 创建日志 |
| `ListOperationLogByClusterID - 按 clusterID 查询` | 按集群查 |
| `ListOperationLogByClusterID - 查询不存在的 clusterID 返回空` | 不存在返回空列表 |
| `ListOperationLogByUserClusterID - 按 clusterID 和用户名查询` | 双重条件 |
| `DeleteOperationLogByTime - 按时间范围删除` | 按时间删除 |

### 5.8 TokenNotify Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateTokenNotify - 创建通知记录` | 创建通知 |
| `GetTokenNotifyByCondition - 按 Token 查询` | 按 Token 查 |
| `GetTokenNotifyByCondition - 查询不存在的 Token 返回空列表` | 不存在返回空 |
| `DeleteTokenNotify - 删除通知记录` | 删除 |

### 5.9 TkeCidr Store（4 方法）

| 测试用例 | 说明 |
|---------|------|
| `SaveTkeCidr - 创建 TkeCidr` | 创建 CIDR |
| `QueryTkeCidr - 按 Vpc 和 Cidr 查询` | 联合查询 |
| `QueryTkeCidr - 查询不存在的记录返回 nil` | 不存在返回 nil |
| `UpdateTkeCidr - 更新 TkeCidr 信息` | 更新状态/集群 |
| `CountTkeCidr - 统计 TkeCidr` | 分组统计 |

### 5.10 User Store（3 方法）

| 测试用例 | 说明 |
|---------|------|
| `CreateUser - 创建用户` | 创建并验证 |
| `GetUserByCondition - 按名称查询用户` | 按名称查 |
| `GetUserByCondition - 查询不存在的用户返回 nil` | 不存在返回 nil |
| `UpdateUser - 更新用户信息` | 更新字段 |

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
| 无 tag（默认） | MySQL 测试 | `go test ./tests/sqlstore/... -v` |
| `dm` | 达梦测试 | `go test -tags=dm ./tests/sqlstore/... -v` |
| `gaussdb` | OpenGauss 测试 | `go test -tags=gaussdb ./tests/sqlstore/... -v` |

### 7.2 添加新测试

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

## 8. 注意事项

### 8.1 GORM AutoMigrate

使用 GORM AutoMigrate 自动创建表结构，无需手动编写 DDL，天然支持多数据库兼容。

### 8.2 全局 GCoreDB

部分 Store 函数依赖全局 `GCoreDB` 变量，测试框架通过 `sqlstore.SetGCoreDB(db)` 设置。

### 8.3 软删除

GORM 默认软删除，`DeleteToken` 等操作实际是更新 `deleted_at` 字段，查询时自动过滤。

### 8.4 OpenSSL 版本要求

**重要**：SM4 加密需要 OpenSSL 1.1.1+。项目依赖 `TencentBlueKing/crypto-golang-sdk`，其中 CGO 代码调用 `EVP_sm4_ctr`，该函数仅在 OpenSSL 1.1.1 及以上版本可用。

- **macOS**：Homebrew 安装的 `openssl@3` 已满足要求
- **Linux CentOS 7**：系统自带的 OpenSSL 1.0.2k 不支持，需要编译安装 OpenSSL 1.1.1+ 到 `/opt/openssl`（见 2.2 节）
- **其他 Linux 发行版**：如果 `openssl version` 输出低于 1.1.1，同样需要升级

### 8.5 调试技巧

- **详细测试步骤**：带上 `-args -ginkgo.v` 参数可以看到每个测试的详细步骤。
- **日志输出**：在测试中使用 `fmt.Fprintln(GinkgoWriter, "your message")` 输出调试信息，输出会与当前测试 Spec 绑定。

## 9. 测试输出示例

### 9.1 汇总模式 (-v)

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
