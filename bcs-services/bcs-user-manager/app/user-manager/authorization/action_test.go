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

import "testing"

func TestNormalizeAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   string
	}{
		{name: "http lowercase", action: " get ", want: "GET"},
		{name: "http mixed case", action: "pAtCh", want: "PATCH"},
		{name: "semantic uppercase", action: " PROJECT_VIEW ", want: ActionProjectView},
		{name: "unknown semantic", action: " Custom_Action ", want: "custom_action"},
		{name: "wildcard", action: " * ", want: "*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeAction(tt.action); got != tt.want {
				t.Fatalf("NormalizeAction(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}

func TestIsKnownAction(t *testing.T) {
	tests := []struct {
		action string
		want   bool
	}{
		{action: "get", want: true},
		{action: " PROJECT_VIEW ", want: true},
		{action: "project_create", want: true},
		{action: "custom_action", want: false},
		{action: "*", want: false},
		{action: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			if got := IsKnownAction(tt.action); got != tt.want {
				t.Fatalf("IsKnownAction(%q) = %t, want %t", tt.action, got, tt.want)
			}
		})
	}
}

func TestResourceTypeForAction(t *testing.T) {
	tests := []struct {
		name             string
		action           string
		wantResourceType string
		wantKnown        bool
	}{
		{name: "project", action: ActionProjectView, wantResourceType: ResourceTypeProject, wantKnown: true},
		{name: "cluster", action: ActionClusterView, wantResourceType: ResourceTypeCluster, wantKnown: true},
		{name: "namespace", action: ActionNamespaceView, wantResourceType: ResourceTypeNamespace, wantKnown: true},
		{name: "global", action: "project_create", wantResourceType: "", wantKnown: true},
		{name: "http action has no semantic resource", action: "GET", wantResourceType: "", wantKnown: false},
		{name: "unknown", action: "custom_action", wantResourceType: "", wantKnown: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceType, known := ResourceTypeForAction(tt.action)
			if resourceType != tt.wantResourceType || known != tt.wantKnown {
				t.Fatalf("ResourceTypeForAction(%q) = (%q, %t), want (%q, %t)",
					tt.action, resourceType, known, tt.wantResourceType, tt.wantKnown)
			}
		})
	}
}
