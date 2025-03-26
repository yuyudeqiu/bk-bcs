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

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/auth"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store/project"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/errorx"
)

// IsMultiTenantEnabled 检查是否启用了多租户模式
func IsMultiTenantEnabled() bool {
	return config.GlobalConf.MultiTenantEnabled
}

// VerifyProject 验证Project是否属于当前租户
func VerifyProject(ctx context.Context, project *project.Project) error {
	// 如果未启用多租户模式，则直接返回nil
	if !IsMultiTenantEnabled() {
		return nil
	}
	if project == nil {
		return errorx.NewReadableErr(errorx.ProjectNotExistsErr, "project is empty")
	}
	tenantID := auth.GetTenantIdFromCtx(ctx)
	// 如果是 default 租户，project 的 tenantID可能为空字符串
	if tenantID == "default" && project.TenantID == "" {
		return nil
	}
	// 如果启用了多租户模式，但资源不属于当前租户，则返回资源不存在错误
	if project.TenantID != tenantID {
		return errorx.NewReadableErr(errorx.ProjectNotExistsErr, "project not found")
	}
	return nil
}
