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
	"errors"
	"testing"
)

type stubBindingReader struct {
	bindings []Binding
	err      error
	calls    int
	subject  string
}

func (s *stubBindingReader) ListBindings(_ context.Context, subject string) ([]Binding, error) {
	s.calls++
	s.subject = subject
	return s.bindings, s.err
}

func TestNoneAuthorizer(t *testing.T) {
	decision, err := (NoneAuthorizer{}).Authorize(context.Background(), Request{Action: "unknown_action"})
	if err != nil {
		t.Fatalf("Authorize() returned error: %v", err)
	}
	if !decision.Allowed {
		t.Fatalf("Authorize() allowed = false, want true")
	}
}

func TestLocalAuthorizerMatching(t *testing.T) { // nolint:funlen
	tests := []struct {
		name     string
		bindings []Binding
		request  Request
		allowed  bool
		reason   string
	}{
		{
			name:     "exact match",
			bindings: []Binding{{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionProjectView}},
			request:  requestFor(ActionProjectView, ResourceTypeProject, "p1", nil),
			allowed:  true,
		},
		{
			name:     "action is normalized",
			bindings: []Binding{{ResourceType: ResourceTypeProject, Resource: "p1", Actions: " PROJECT_VIEW "}},
			request:  requestFor(" Project_View ", ResourceTypeProject, "p1", nil),
			allowed:  true,
		},
		{
			name:     "http action is normalized",
			bindings: []Binding{{ResourceType: ResourceTypeCluster, Resource: "c1", Actions: "get"}},
			request:  requestFor(" GET ", ResourceTypeCluster, "c1", nil),
			allowed:  true,
		},
		{
			name:     "action wildcard",
			bindings: []Binding{{ResourceType: ResourceTypeProject, Resource: "p1", Actions: "*"}},
			request:  requestFor("project_edit", ResourceTypeProject, "p1", nil),
			allowed:  true,
		},
		{
			name:     "resource type wildcard",
			bindings: []Binding{{ResourceType: "*", Resource: "c1", Actions: ActionClusterView}},
			request:  requestFor(ActionClusterView, ResourceTypeCluster, "c1", nil),
			allowed:  true,
		},
		{
			name:     "resource wildcard",
			bindings: []Binding{{ResourceType: ResourceTypeCluster, Resource: "*", Actions: ActionClusterView}},
			request:  requestFor(ActionClusterView, ResourceTypeCluster, "c1", nil),
			allowed:  true,
		},
		{
			name:     "one of multiple actions matches",
			bindings: []Binding{{ResourceType: ResourceTypeCluster, Resource: "c1", Actions: "cluster_use, cluster_view"}},
			request:  requestFor(ActionClusterView, ResourceTypeCluster, "c1", nil),
			allowed:  true,
		},
		{
			name:     "wrong resource is denied",
			bindings: []Binding{{ResourceType: ResourceTypeProject, Resource: "p2", Actions: ActionProjectView}},
			request:  requestFor(ActionProjectView, ResourceTypeProject, "p1", nil),
			reason:   "no matching role binding",
		},
		{
			name:     "wrong action is denied",
			bindings: []Binding{{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionProjectView}},
			request:  requestFor("project_edit", ResourceTypeProject, "p1", nil),
			reason:   "no matching role binding",
		},
		{
			name:     "unknown action is denied despite wildcard",
			bindings: []Binding{{ResourceType: "*", Resource: "*", Actions: "*"}},
			request:  requestFor("custom_action", ResourceTypeProject, "p1", nil),
			reason:   "unknown action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &stubBindingReader{bindings: tt.bindings}
			decision, err := NewLocalAuthorizer(reader).Authorize(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("Authorize() returned error: %v", err)
			}
			if decision.Allowed != tt.allowed || decision.Reason != tt.reason {
				t.Fatalf("Authorize() = %+v, want allowed=%t reason=%q", decision, tt.allowed, tt.reason)
			}
		})
	}
}

func TestLocalAuthorizerResourceInheritance(t *testing.T) { // nolint:funlen
	tests := []struct {
		name    string
		binding Binding
		request Request
		allowed bool
	}{
		{
			name:    "cluster inherits project binding",
			binding: Binding{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionClusterView},
			request: requestFor(ActionClusterView, ResourceTypeCluster, "c1", map[string]string{
				AttributeProjectID: "p1",
			}),
			allowed: true,
		},
		{
			name:    "cluster without project context does not inherit",
			binding: Binding{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionClusterView},
			request: requestFor(ActionClusterView, ResourceTypeCluster, "c1", nil),
		},
		{
			name:    "namespace inherits cluster binding",
			binding: Binding{ResourceType: ResourceTypeCluster, Resource: "c1", Actions: ActionNamespaceView},
			request: requestFor(ActionNamespaceView, ResourceTypeNamespace, "ns1", map[string]string{
				AttributeClusterID: "c1",
				AttributeProjectID: "p1",
			}),
			allowed: true,
		},
		{
			name:    "namespace inherits project binding",
			binding: Binding{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionNamespaceView},
			request: requestFor(ActionNamespaceView, ResourceTypeNamespace, "ns1", map[string]string{
				AttributeClusterID: "c1",
				AttributeProjectID: "p1",
			}),
			allowed: true,
		},
		{
			name:    "namespace scoped resource inherits cluster binding",
			binding: Binding{ResourceType: ResourceTypeCluster, Resource: "c1", Actions: "namespace_scoped_view"},
			request: requestFor("namespace_scoped_view", ResourceTypeNamespaceScoped, "pods", map[string]string{
				AttributeClusterID: "c1",
			}),
			allowed: true,
		},
		{
			name:    "parent does not inherit child binding",
			binding: Binding{ResourceType: ResourceTypeCluster, Resource: "c1", Actions: ActionProjectView},
			request: requestFor(ActionProjectView, ResourceTypeProject, "p1", map[string]string{
				AttributeClusterID: "c1",
			}),
		},
		{
			name:    "action must match when resource is inherited",
			binding: Binding{ResourceType: ResourceTypeProject, Resource: "p1", Actions: ActionProjectView},
			request: requestFor(ActionClusterView, ResourceTypeCluster, "c1", map[string]string{
				AttributeProjectID: "p1",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &stubBindingReader{bindings: []Binding{tt.binding}}
			decision, err := NewLocalAuthorizer(reader).Authorize(context.Background(), tt.request)
			if err != nil {
				t.Fatalf("Authorize() returned error: %v", err)
			}
			if decision.Allowed != tt.allowed {
				t.Fatalf("Authorize() allowed = %t, want %t", decision.Allowed, tt.allowed)
			}
		})
	}
}

func TestLocalAuthorizerSuperuserSkipsBindingLookup(t *testing.T) {
	reader := &stubBindingReader{err: errors.New("must not be returned")}
	decision, err := NewLocalAuthorizer(reader).Authorize(context.Background(), Request{
		Subject:   "admin",
		Superuser: true,
		Action:    "unknown_action",
	})
	if err != nil {
		t.Fatalf("Authorize() returned error: %v", err)
	}
	if !decision.Allowed || decision.Reason != "superuser" {
		t.Fatalf("Authorize() = %+v, want superuser allow", decision)
	}
	if reader.calls != 0 {
		t.Fatalf("ListBindings() calls = %d, want 0", reader.calls)
	}
}

func TestLocalAuthorizerPassesSubjectAndReturnsReaderError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	reader := &stubBindingReader{err: wantErr}
	decision, err := NewLocalAuthorizer(reader).Authorize(context.Background(), requestFor(
		ActionProjectView, ResourceTypeProject, "p1", nil))
	if !errors.Is(err, wantErr) {
		t.Fatalf("Authorize() error = %v, want %v", err, wantErr)
	}
	if decision != (Decision{}) {
		t.Fatalf("Authorize() decision = %+v, want zero value", decision)
	}
	if reader.subject != "alice" {
		t.Fatalf("ListBindings() subject = %q, want %q", reader.subject, "alice")
	}
}

func requestFor(action, resourceType, resourceID string, attributes map[string]string) Request {
	return Request{
		Subject: "alice",
		Action:  action,
		Resource: Resource{
			Type:       resourceType,
			ID:         resourceID,
			Attributes: attributes,
		},
	}
}
