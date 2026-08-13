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
	resourceTypeProject      = authorization.ResourceTypeProject
	resourceTypeCluster      = authorization.ResourceTypeCluster
	resourceTypeNamespace    = authorization.ResourceTypeNamespace
	resourceTypeTemplateSet  = authorization.ResourceTypeTemplateSet
	resourceTypeCloudAccount = authorization.ResourceTypeCloudAccount
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
			authorization.AttributeProjectID: permCtx.ProjectID,
			authorization.AttributeClusterID: permCtx.ClusterID,
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
func GetResourceTypeFromAction(action string) string {
	resourceType, _ := authorization.ResourceTypeForAction(action)
	return resourceType
}
