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

package auth

import (
	"encoding/json"

	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/cloudaccount"
	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/cluster"
	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/namespace"
	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/project"
	"github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/templateset"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/authorization"
)

// PermCtx perm context
type PermCtx struct {
	ResourceType string      `json:"resource_type"`
	ProjectID    string      `json:"project_id"`
	ClusterID    string      `json:"cluster_id"`
	Namespace    string      `json:"name"`
	TemplateID   json.Number `json:"template_id"`
	AccountID    string      `json:"account_id"`
}

// ResourceFromPermCtx converts the compatibility API context into an internal resource.
func ResourceFromPermCtx(permCtx *PermCtx) authorization.Resource {
	if permCtx == nil {
		return authorization.Resource{}
	}
	resource := authorization.Resource{
		Type: permCtx.ResourceType,
		Attributes: map[string]string{
			"project_id": permCtx.ProjectID,
			"cluster_id": permCtx.ClusterID,
		},
	}
	switch permCtx.ResourceType {
	case string(project.SysProject):
		resource.ID = permCtx.ProjectID
	case string(cluster.SysCluster):
		resource.ID = permCtx.ClusterID
	case string(namespace.SysNamespace):
		resource.ID = permCtx.Namespace
	case string(templateset.SysTemplateSet):
		resource.ID = permCtx.TemplateID.String()
	case string(cloudaccount.SysCloudAccount):
		resource.ID = permCtx.AccountID
	}
	return resource
}

// GetResourceTypeFromAction get resource type from action
// NOCC:CCN_threshold(工具误报:),golint/fnsize(设计如此:)
func GetResourceTypeFromAction(action string) string { // nolint
	switch action {
	case project.ProjectCreate.String():
		return ""
	case project.ProjectView.String():
		return string(project.SysProject)
	case project.ProjectEdit.String():
		return string(project.SysProject)
	case project.ProjectDelete.String():
		return string(project.SysProject)
	case cluster.ClusterCreate.String():
		return string(project.SysProject)
	case cluster.ClusterView.String():
		return string(cluster.SysCluster)
	case cluster.ClusterManage.String():
		return string(cluster.SysCluster)
	case cluster.ClusterDelete.String():
		return string(cluster.SysCluster)
	case cluster.ClusterUse.String():
		return string(cluster.SysCluster)
	case namespace.NameSpaceCreate.String():
		return string(cluster.SysCluster)
	case namespace.NameSpaceView.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceUpdate.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceDelete.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceList.String():
		return string(cluster.SysCluster)
	case cluster.ClusterScopedCreate.String():
		return string(cluster.SysCluster)
	case cluster.ClusterScopedView.String():
		return string(cluster.SysCluster)
	case cluster.ClusterScopedUpdate.String():
		return string(cluster.SysCluster)
	case cluster.ClusterScopedDelete.String():
		return string(cluster.SysCluster)
	case namespace.NameSpaceScopedCreate.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceScopedView.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceScopedUpdate.String():
		return string(namespace.SysNamespace)
	case namespace.NameSpaceScopedDelete.String():
		return string(namespace.SysNamespace)
	case templateset.TemplateSetCreate.String():
		return string(project.SysProject)
	case templateset.TemplateSetView.String():
		return string(templateset.SysTemplateSet)
	case templateset.TemplateSetCopy.String():
		return string(templateset.SysTemplateSet)
	case templateset.TemplateSetUpdate.String():
		return string(templateset.SysTemplateSet)
	case templateset.TemplateSetDelete.String():
		return string(templateset.SysTemplateSet)
	case templateset.TemplateSetInstantiate.String():
		return string(templateset.SysTemplateSet)
	case cloudaccount.AccountCreate.String():
		return string(project.SysProject)
	case cloudaccount.AccountManage.String():
		return string(cloudaccount.SysCloudAccount)
	case cloudaccount.AccountUse.String():
		return string(cloudaccount.SysCloudAccount)
	default:
		return ""
	}
}
