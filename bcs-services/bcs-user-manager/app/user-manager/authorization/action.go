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

package authorization

import (
	"net/http"
	"strings"
)

const (
	// ActionHTTPGet is retained for the V2 HTTP-method permission API.
	ActionHTTPGet = http.MethodGet

	// Frequently used read actions are exported for built-in roles and filters.
	ActionProjectView   = "project_view"
	ActionClusterView   = "cluster_view"
	ActionClusterUse    = "cluster_use"
	ActionNamespaceView = "namespace_view"
	ActionNamespaceList = "namespace_list"
)

var httpActions = map[string]struct{}{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodPut:    {},
	http.MethodPatch:  {},
	http.MethodDelete: {},
}

// semanticActions is the compatibility action catalog. An empty resource type
// means that the action is not scoped to an existing resource.
var semanticActions = map[string]string{
	"project_create":  "",
	ActionProjectView: ResourceTypeProject,
	"project_edit":    ResourceTypeProject,
	"project_delete":  ResourceTypeProject,

	"cluster_create":       ResourceTypeProject,
	"namespace_create":     ResourceTypeProject,
	"templateset_create":   ResourceTypeProject,
	"cloud_account_create": ResourceTypeProject,

	ActionClusterView:       ResourceTypeCluster,
	"cluster_manage":        ResourceTypeCluster,
	"cluster_delete":        ResourceTypeCluster,
	ActionClusterUse:        ResourceTypeCluster,
	ActionNamespaceList:     ResourceTypeCluster,
	"cluster_scoped_create": ResourceTypeCluster,
	"cluster_scoped_view":   ResourceTypeCluster,
	"cluster_scoped_update": ResourceTypeCluster,
	"cluster_scoped_delete": ResourceTypeCluster,

	ActionNamespaceView:       ResourceTypeNamespace,
	"namespace_update":        ResourceTypeNamespace,
	"namespace_delete":        ResourceTypeNamespace,
	"namespace_scoped_create": ResourceTypeNamespace,
	"namespace_scoped_view":   ResourceTypeNamespace,
	"namespace_scoped_update": ResourceTypeNamespace,
	"namespace_scoped_delete": ResourceTypeNamespace,

	"templateset_view":        ResourceTypeTemplateSet,
	"templateset_copy":        ResourceTypeTemplateSet,
	"templateset_update":      ResourceTypeTemplateSet,
	"templateset_delete":      ResourceTypeTemplateSet,
	"templateset_instantiate": ResourceTypeTemplateSet,

	"cloud_account_manage": ResourceTypeCloudAccount,
	"cloud_account_use":    ResourceTypeCloudAccount,
}

// NormalizeAction produces the representation stored and evaluated locally.
func NormalizeAction(action string) string {
	action = strings.TrimSpace(action)
	upper := strings.ToUpper(action)
	if _, ok := httpActions[upper]; ok {
		return upper
	}
	return strings.ToLower(action)
}

// IsKnownAction reports whether an action is part of the local authorization contract.
func IsKnownAction(action string) bool {
	action = NormalizeAction(action)
	if _, ok := httpActions[action]; ok {
		return true
	}
	_, ok := semanticActions[action]
	return ok
}

// ResourceTypeForAction returns the resource type associated with a semantic action.
func ResourceTypeForAction(action string) (string, bool) {
	resourceType, ok := semanticActions[NormalizeAction(action)]
	return resourceType, ok
}
