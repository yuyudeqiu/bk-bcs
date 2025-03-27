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

package project

import (
	"context"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/auth"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/errorx"
	proto "github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/proto/bcsproject"
)

// DeleteAction action for delete project
type DeleteAction struct {
	ctx   context.Context
	model store.ProjectModel
	req   *proto.DeleteProjectRequest
}

// NewDeleteAction delete project action
func NewDeleteAction(model store.ProjectModel) *DeleteAction {
	return &DeleteAction{
		model: model,
	}
}

// Do delete project
func (da *DeleteAction) Do(ctx context.Context, req *proto.DeleteProjectRequest) error {
	da.ctx = ctx
	da.req = req

	// 不能删除别的租户的 Project
	project, err := da.model.GetProject(ctx, req.GetProjectID())
	if err != nil {
		return errorx.NewDBErr(err.Error())
	}
	// 未找到，或者项目属于别的租户，都不能删除
	// TODO 通过一个数据库操作完成，不用先查
	if project == nil || (config.GlobalConf.MultiTenantEnabled && project.TenantID != auth.GetTenantFromCtx(ctx)) {
		return errorx.NewReadableErr(errorx.ProjectNotExistsErr, "project not found")
	}

	if err := da.model.DeleteProject(ctx, req.ProjectID); err != nil {
		return errorx.NewDBErr(err.Error())
	}

	return nil
}
