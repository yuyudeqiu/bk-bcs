/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package tenant xxx
package tenant

import (
	"context"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/common/constant"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/common/headerkey"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/config"
)

// IsMultiTenantEnabled 检查是否启用了多租户模式
func IsMultiTenantEnabled() bool {
	return config.GlobalConf.MultiTenantEnabled
}

// GetTenantIdFromCtx get tenantId from context
func GetTenantIdFromCtx(ctx context.Context) string {
	tenantId := ""
	if id, ok := ctx.Value(headerkey.TenantIdKey).(string); ok {
		tenantId = id
	}
	if tenantId == "" {
		tenantId = constant.DefaultTenantId
	}
	return tenantId
}
