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

// Package project xxx
package project

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/auth"
	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/middleware"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/component/bcscc"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/logging"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store"
	pm "github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store/project"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/errorx"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/stringx"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/tenant"
	proto "github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/proto/bcsproject"
)

// CreateAction action for create project
type CreateAction struct {
	ctx   context.Context
	model store.ProjectModel
	req   *proto.CreateProjectRequest
}

// NewCreateAction new create project action
func NewCreateAction(model store.ProjectModel) *CreateAction {
	return &CreateAction{
		model: model,
	}
}

// Do create project request
func (ca *CreateAction) Do(ctx context.Context, req *proto.CreateProjectRequest) (*pm.Project, error) {
	ca.ctx = ctx
	ca.req = req

	if err := ca.validate(); err != nil {
		return nil, errorx.NewReadableErr(errorx.ParamErr, err.Error())
	}

	// 如果有传递项目ID，则以传递的为准，否则动态生成32位的字符串作为项目ID
	if req.ProjectID == "" {
		ca.req.ProjectID = stringx.GenUUID()
	}

	if err := ca.createProject(); err != nil {
		return nil, errorx.NewDBErr(err.Error())
	}

	p, err := ca.model.GetProject(ca.ctx, ca.req.ProjectID)
	if err != nil {
		return nil, errorx.NewDBErr(err.Error())
	}
	// 向 bcs cc 写入数据
	go func() {
		if err := bcscc.CreateProject(p); err != nil {
			logging.Error("[ALARM-CC-PROJECT] create project %s/%s in paas-cc failed, err: %s",
				p.ProjectID, p.ProjectCode, err.Error())
		}

	}()
	// 返回项目信息
	return p, nil
}

func (ca *CreateAction) createProject() error {
	p := &pm.Project{
		ProjectID:   ca.req.ProjectID,
		Name:        ca.req.Name,
		ProjectCode: ca.req.ProjectCode,
		ProjectType: ca.req.ProjectType,
		UseBKRes:    ca.req.UseBKRes,
		Description: ca.req.Description,
		IsOffline:   ca.req.IsOffline,
		Kind:        ca.req.Kind,
		BusinessID:  ca.req.BusinessID,
		DeployType:  ca.req.DeployType,
		BGID:        ca.req.BGID,
		BGName:      ca.req.BGName,
		DeptID:      ca.req.DeptID,
		DeptName:    ca.req.DeptName,
		CenterID:    ca.req.CenterID,
		CenterName:  ca.req.CenterName,
		IsSecret:    ca.req.IsSecret,
		Labels:      ca.req.Labels,
		Annotations: ca.req.Annotations,
		CreateTime:  time.Now().Format(time.RFC3339),
		UpdateTime:  time.Now().Format(time.RFC3339),
	}
	// 从 context 中获取 username
	if authUser, err := middleware.GetUserFromContext(ca.ctx); err == nil {
		p.Creator = authUser.GetUsername()
		p.Managers = authUser.GetUsername()
		p.TenantID = authUser.GetTanantId() // 单租户模式下该字段为 default
	}

	// TODO 使用 user-manager 接口获取 用户当前所属的租户 & displayName & bk_username
	// 单租户 使用 projectCode 且 projectCode = tenantProjectCode，保持原 p.Project = ca.req.Project 即可
	// 多租户 使用 tenantProjectCode 且后台转换至 projectCode。
	// projectCode 全局唯一，所以需要拼接 租户信息+租户视角下的 tenantProjectCode，并且将拼接后的结果写到 ProjectCode
	if tenant.IsMultiTenantEnabled() {
		p.TenantProjectCode = p.ProjectCode
		p.ProjectCode = ca.generateProjectCode(p.TenantID, p.TenantProjectCode)
	}

	return ca.model.CreateProject(ca.ctx, p)
}

func (ca *CreateAction) validate() error {
	// check projectID、projectCode、name
	projectID, projectCode, name := ca.req.ProjectID, ca.req.ProjectCode, ca.req.Name
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("name cannot contains only spaces")
	}
	// 先检查全局唯一的 ProjectID是否已经被使用了
	p, err := ca.model.GetProject(ca.ctx, projectID)
	if err != nil {
		return err
	}
	if p != nil {
		return fmt.Errorf("projectID: %s is already exists", projectID)
	}

	if tenant.IsMultiTenantEnabled() {
		// 多租户 TenantID+TenantProjectCode 全局唯一，其他条件均为 Or
		p, _ = ca.model.GetTenantProjectByField(ca.ctx, &pm.ProjectField{TenantID: auth.GetTenantIdFromCtx(ca.ctx),
			TenantProjectCode: projectCode, Name: name})
		if p.TenantProjectCode == projectCode {
			return fmt.Errorf("projectCode: %s is already exists", projectCode)
		}
		if p.Name == name {
			return fmt.Errorf("name: %s is already exists", name)
		}
		return nil
	}
	// 单租户情况保持原查询方式
	if p, _ = ca.model.GetProjectByField(ca.ctx, &pm.ProjectField{ProjectID: projectID, ProjectCode: projectCode,
		Name: name}); p != nil {
		if p.ProjectID == projectID {
			return fmt.Errorf("projectID: %s is already exists", projectID)
		}
		if p.ProjectCode == projectCode {
			return fmt.Errorf("projectCode: %s is already exists", projectCode)
		}
		if p.Name == name {
			return fmt.Errorf("name: %s is already exists", name)
		}
	}
	return nil
}

func (ca *CreateAction) generateProjectCode(tenantID, tenantProjectCode string) string {
	// TODO 系统生成，english_name=xxxx-$｛tenant_english_name｝ 前缀为租户 ID，分隔符为中划线（-）
	//  这里可能需要一些 user-manager 提供的信息来拼接？
	return fmt.Sprintf("%s-%s", tenantID, tenantProjectCode)
}
