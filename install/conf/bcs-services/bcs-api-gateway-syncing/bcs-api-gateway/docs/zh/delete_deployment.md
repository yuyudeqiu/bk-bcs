### 描述

删除 Deployment

### 路径参数

| 名称        | 位置   | 类型     | 必选 | 说明            |
|-----------|------|--------|----|---------------|
| projectID | path | string | 是  | 项目 ID         |
| clusterID | path | string | 是  | 集群 ID         |
| namespace | path | string | 是  | 命名空间          |
| name      | path | string | 是  | Deployment 名称 |

### 调用示例

```sh
curl -X DELETE \
-H 'Cookie: bk_token=xxx' \
-H 'User-Agent: xxx' \
'http://bcs-api.bkdomain/bcsapi/v4/clusterresources/v1/projects/{projectID}/clusters/{clusterID}/namespaces/{namespace}/workloads/deployments/{name}'
```

### 响应示例

```json
{
  "code": 0,
  "message": "string",
  "requestID": "string",
  "data": {},
  "webAnnotations": {}
}
```