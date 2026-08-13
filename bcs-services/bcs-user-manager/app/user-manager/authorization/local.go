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
	"context"
	"strings"
)

// LocalAuthorizer evaluates role bindings stored by user-manager.
type LocalAuthorizer struct {
	bindings BindingReader
}

// NewLocalAuthorizer creates a database-backed authorizer.
func NewLocalAuthorizer(bindings BindingReader) *LocalAuthorizer {
	return &LocalAuthorizer{bindings: bindings}
}

// Authorize allows superusers and subjects with a matching resource/action binding.
func (a *LocalAuthorizer) Authorize(ctx context.Context, request Request) (Decision, error) {
	if request.Superuser {
		return Decision{Allowed: true, Reason: "superuser"}, nil
	}
	action := NormalizeAction(request.Action)
	if !IsKnownAction(action) {
		return Decision{Allowed: false, Reason: "unknown action"}, nil
	}

	bindings, err := a.bindings.ListBindings(ctx, request.Subject)
	if err != nil {
		return Decision{}, err
	}
	for _, resource := range resourceCandidates(request.Resource) {
		for _, binding := range bindings {
			if !matches(binding.ResourceType, resource.Type) ||
				!matches(binding.Resource, resource.ID) {
				continue
			}
			for _, rule := range strings.Split(binding.Actions, ",") {
				if matches(NormalizeAction(rule), action) {
					return Decision{Allowed: true}, nil
				}
			}
		}
	}
	return Decision{Allowed: false, Reason: "no matching role binding"}, nil
}

func matches(rule, value string) bool {
	return rule == "*" || rule == value
}

// resourceCandidates returns the requested resource followed by its known
// ancestors. Missing hierarchy attributes are ignored instead of guessed.
func resourceCandidates(resource Resource) []Resource {
	resources := []Resource{{Type: resource.Type, ID: resource.ID}}
	appendResource := func(resourceType, id string) {
		if id == "" {
			return
		}
		resources = append(resources, Resource{Type: resourceType, ID: id})
	}

	switch resource.Type {
	case ResourceTypeCluster:
		appendResource(ResourceTypeProject, resource.Attributes[AttributeProjectID])
	case ResourceTypeNamespace, ResourceTypeNamespaceScoped:
		appendResource(ResourceTypeCluster, resource.Attributes[AttributeClusterID])
		appendResource(ResourceTypeProject, resource.Attributes[AttributeProjectID])
	}
	return resources
}
