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

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/authorization"
)

const (
	resourceTypeProject      = "project"
	resourceTypeCluster      = "cluster"
	resourceTypeNamespace    = "namespace"
	resourceTypeTemplateSet  = "templateset"
	resourceTypeCloudAccount = "cloud_account"
)

// PermCtx is the resource context accepted by the compatibility permission API.
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
	case resourceTypeProject:
		resource.ID = permCtx.ProjectID
	case resourceTypeCluster:
		resource.ID = permCtx.ClusterID
	case resourceTypeNamespace:
		resource.ID = permCtx.Namespace
	case resourceTypeTemplateSet:
		resource.ID = permCtx.TemplateID.String()
	case resourceTypeCloudAccount:
		resource.ID = permCtx.AccountID
	}
	return resource
}

// GetResourceTypeFromAction returns the resource protected by a compatibility action.
func GetResourceTypeFromAction(action string) string { // nolint:cyclop
	switch action {
	case "project_create":
		return ""
	case "project_view", "project_edit", "project_delete", "cluster_create", "namespace_create",
		"templateset_create", "cloud_account_create":
		return resourceTypeProject
	case "cluster_view", "cluster_manage", "cluster_delete", "cluster_use", "namespace_list",
		"cluster_scoped_create", "cluster_scoped_view", "cluster_scoped_update", "cluster_scoped_delete":
		return resourceTypeCluster
	case "namespace_view", "namespace_update", "namespace_delete", "namespace_scoped_create",
		"namespace_scoped_view", "namespace_scoped_update", "namespace_scoped_delete":
		return resourceTypeNamespace
	case "templateset_view", "templateset_copy", "templateset_update", "templateset_delete",
		"templateset_instantiate":
		return resourceTypeTemplateSet
	case "cloud_account_manage", "cloud_account_use":
		return resourceTypeCloudAccount
	default:
		return ""
	}
}
