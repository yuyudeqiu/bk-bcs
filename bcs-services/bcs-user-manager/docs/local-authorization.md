# 本地权限改造说明

本文档记录 `bcs-user-manager` 权限改造的当前状态，方便后续贡献者或 Agent 在不了解原蓝鲸 IAM
实现细节的情况下继续开发。

## 改造目标

让 `bcs-user-manager` 在不依赖蓝鲸 IAM 的情况下运行，同时保留一个足够轻量的权限边界，支持以下模式：

- `none`：验证用户身份后，放行所有权限请求。
- `local`：使用 user-manager 自己数据库中的角色和资源绑定进行鉴权。

本项目需要的是轻量 RBAC，不是重新实现一套蓝鲸 IAM。

## 已完成的改造

目前四个阶段对应以下提交：

- `02acebe48 refactor(user-manager): introduce authorization boundary`
- `8c0bf2fa9 refactor(user-manager): remove blueking iam integration`
- `3852e0e82 feat(user-manager): add local authorization mode`
- `b09f9497d feat(user-manager): strengthen local authorization`

目前已经删除：

- 蓝鲸 IAM 客户端及其初始化逻辑
- IAM 资源 Provider 回调及其实现
- IAM migration 文件和启动迁移逻辑
- `iam_config` 和 `permission_switch`
- `iam-go-sdk`、`bcs-services/pkg`、`golang-migrate` 等相关直接依赖
- 原有 IAM/Kubernetes 权限转换和校验代码

提交 `b09f9497d` 在本地鉴权基础上继续完成：

- 将角色绑定从 token 记录 ID 改为稳定的 subject 名称
- 增加角色绑定联合唯一索引，并让数据库 schema 迁移错误阻止服务启动
- 增加 project → cluster → namespace 的单向资源继承
- 建立本地 action 目录，对未知 action 默认拒绝
- 将 TKE CIDR 分配使用的锁从已删除的权限缓存中解耦，迁回 TKE 包内部

后续已经为纯 `authorization` 包补充单元测试，覆盖 `none`、factory、action 目录、通配符、超级用户、subject
透传、资源继承、未知 action 拒绝和存储错误传播。

原 IAM 资源回调接口 `/v1/iam-provider/resources` 已经删除。

以下两个地址为了兼容可能存在的 BCS 前端或客户端调用而暂时保留：

- `POST /v1/iam/user_perms`
- `POST /v1/iam/user_perms/actions/{action_id}`

它们不是 IAM 回调。当前已经改为调用内部 `Authorizer`，路径中的 `iam` 只是历史命名。

权限实现位于：

```text
app/user-manager/authorization
```

## 当前本地权限模型

`local` 模式复用以下数据表：

- `bcs_users`
- `bcs_roles`
- `bcs_user_resource_roles`

角色中的 action 暂时使用逗号分隔字符串。资源绑定使用稳定的逻辑名称 `subject`，把一个用户或 client、一个角色
和特定资源关联起来；它不再关联 `bcs_users.id`，避免 token 记录被删除或重建时权限一起丢失。

当前匹配规则：

- action 支持精确匹配或 `*`
- resource type 支持精确匹配或 `*`
- resource ID 支持精确匹配或 `*`
- 用户存在多个绑定时，任意一个匹配即可放行
- admin 和 saas 用户始终放行
- 没有匹配绑定时默认拒绝
- 暂不支持显式 deny

绑定表使用 `(subject, resource_type, resource, role_id)` 联合唯一索引，避免重复授权。该索引的左前缀也覆盖
当前按用户鉴权、按用户和资源类型查询绑定的主要访问方式。`resource` 为非空字段；不限定具体资源时使用 `*`，
而不是 `NULL`。

全新数据库会创建两个内置角色：

- `manager`：`*`
- `viewer`：`GET,project_view,cluster_view,cluster_use,namespace_view,namespace_list`

当前每次鉴权直接查询 MySQL，不再使用原来每 60 秒刷新的全局权限缓存，也删除了启动时等待缓存的
`time.Sleep`。除非后续压测证明确有必要，否则不需要过早加入缓存。

### 资源继承

本地鉴权支持单向的父资源权限继承：

```text
project → cluster → namespace / namespace_scoped
```

- cluster 请求先匹配自身绑定，再使用请求携带的 `project_id` 匹配项目绑定。
- namespace 和 namespace_scoped 请求先匹配自身，再依次使用 `cluster_id`、`project_id` 匹配父级绑定。
- 继承只改变参与匹配的资源，action 仍然必须匹配。例如项目角色包含 `cluster_view` 时才能查看其下集群，只有
  `project_view` 不会自动获得 `cluster_view`。
- 请求没有携带父资源 ID 时，只进行已有层级的匹配；Authorizer 不查询其他 BCS 服务，也不推测资源归属。
- V1 verify 请求没有父资源上下文，因此保持精确匹配。V2 verify 和携带 `PermCtx` 的兼容接口支持继承。
- 不支持子资源权限向父资源反向继承。

### 权限主体约定

本地权限使用逻辑名称作为 subject，不使用 `bcs_users.id` 作为对外的 subject ID：

- 普通用户 token 使用用户名。
- 临时 token 使用它所代表的 `Username`，权限随目标用户继承；`CreatedBy` 只保留为签发审计信息。
- client token 使用 client name。

`bcs_users` 中的一行可能只是某个用户的一条 token 记录，记录 ID 并不是稳定的逻辑身份 ID。绑定表因此直接保存
subject 名称，不依赖 token 记录的生命周期。用户名和 client name 暂时共享同一个命名空间；新接口需要阻止两类
主体创建同名记录，暂不为此引入更复杂的 subject 类型。

### Action 目录

本地权限明确登记两类 action：

- 业务 action，例如 `project_view`、`cluster_manage`、`namespace_view`。每个业务 action 同时登记其资源类型。
- V2 verify 兼容使用的 HTTP action：`GET`、`POST`、`PUT`、`PATCH`、`DELETE`。

两类 action 不做自动互译，因为仅凭 HTTP 方法无法判断业务语义。鉴权前会去除首尾空格，将 HTTP action 统一为
大写、业务 action 统一为小写。`local` 模式默认拒绝未登记 action；即使角色 action 为 `*` 也不会放行未知
action。新增业务能力时需要显式更新 `authorization/action.go`，使权限面的变化可以被审查。`none` 模式仍然放行
所有已经完成身份认证的请求。

## 配置方式

启用本地权限：

```json
{
  "authorization": {
    "mode": "local"
  }
}
```

容器模板使用环境变量 `authorizationMode`，默认值为 `none`。

未知 mode 会导致 user-manager 启动失败，避免因为配置拼写错误而意外放行。

## 当前权限管理接口

暂时复用原有接口管理角色绑定：

- `POST /v1/permissions`：为用户授予资源角色
- `GET /v1/permissions`：查询用户在某类资源上的角色绑定
- `DELETE /v1/permissions`：撤销资源角色
- `GET /v1/permissions/verify`：V1 兼容权限验证接口
- `GET /v2/permissions/verify`：V2 兼容权限验证接口

V1 和 V2 验证接口都会先验证 token，再调用当前配置的 `Authorizer`。

## 部署前提

该分支只面向全新部署和空数据库，不支持从已有蓝鲸 BCS user-manager 数据库升级，也不需要保留历史
IAM 数据或历史角色配置。

数据库初始化只需要满足服务正常重启时的幂等性。后续可以直接按照新的权限模型调整表结构，不需要增加
旧数据迁移和兼容逻辑。

## 尚未完成的工作

建议后续按独立、可审查的提交继续：

1. 设计更清晰的本地角色和绑定管理 API。在确认所有调用方之前先保留旧接口，并阻止用户名与 client name 同名。
2. 后续增加 `/v1/authorization/...` 新路径，并逐步废弃两个历史 `/v1/iam/...` 路径。
3. 决定拒绝结果是否需要结构化原因和审计记录。本地版本不再提供 IAM 权限申请地址。
4. 为 `AuthorizationStore` 的 SQL 查询和 HTTP 兼容接口补充测试，并在依赖继续精简后运行完整测试集。

## 当前验证状态

在提交 `b09f9497d` 前，已经在 `bcs-services/bcs-user-manager` 目录执行：

```text
go build .
```

构建成功。构建过程中发现原 TKE CIDR 逻辑仍引用已随权限缓存删除的 `permission.Mutex`，现已改为 TKE 包内部
的 `sync.Mutex`。构建生成的本地二进制已经清理，没有进入提交。

纯鉴权包已经执行：

```text
go test ./app/user-manager/authorization
go test -cover ./app/user-manager/authorization
```

测试通过，语句覆盖率为 `100.0%`。当前尚未运行数据库存储测试、HTTP 接口测试和整个 user-manager 的完整测试集。

## 当前明确不做的内容

- 不恢复 IAM 资源注册和权限申请流程
- 不实现显式 deny 规则
- 不设计策略语言
- 不在 user-manager 本地权限稳定前扩展到所有 BCS 服务

后续应继续保持 `Authorizer` 接口足够小。新的本地能力应放在接口实现后面，不要再把数据库类型或蓝鲸专属
类型泄漏到 HTTP Filter 和业务 Handler 中。
