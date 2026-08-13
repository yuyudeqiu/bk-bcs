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

// Package authorization defines the user-manager authorization boundary.
package authorization

import "context"

const (
	// ModeNone disables authorization checks after the caller has been authenticated.
	ModeNone = "none"
	// ModeLocal evaluates role bindings stored in the user-manager database.
	ModeLocal = "local"

	// ResourceTypeProject is the root resource in the local BCS hierarchy.
	ResourceTypeProject = "project"
	// ResourceTypeCluster is a child of a project.
	ResourceTypeCluster = "cluster"
	// ResourceTypeNamespace is a child of a cluster.
	ResourceTypeNamespace = "namespace"
	// ResourceTypeNamespaceScoped represents resources contained by a namespace.
	ResourceTypeNamespaceScoped = "namespace_scoped"
	// ResourceTypeTemplateSet represents a BCS template set.
	ResourceTypeTemplateSet = "templateset"
	// ResourceTypeCloudAccount represents a cloud account.
	ResourceTypeCloudAccount = "cloud_account"

	// AttributeProjectID identifies the project containing a resource.
	AttributeProjectID = "project_id"
	// AttributeClusterID identifies the cluster containing a resource.
	AttributeClusterID = "cluster_id"
)

// Resource identifies the object affected by an action.
type Resource struct {
	Type       string
	ID         string
	Attributes map[string]string
}

// Request contains the information required to make an authorization decision.
type Request struct {
	// Subject is the effective username or client name. Temporary tokens use the
	// username they represent instead of the client that created the token.
	Subject   string
	Superuser bool
	Action    string
	Resource  Resource
}

// Decision is the result returned by an Authorizer.
type Decision struct {
	Allowed bool
	Reason  string
}

// Authorizer is implemented by each authorization mode.
type Authorizer interface {
	Authorize(ctx context.Context, request Request) (Decision, error)
}

// Binding is a resolved role binding used by the local authorizer.
type Binding struct {
	ResourceType string
	Resource     string
	Actions      string
}

// BindingReader loads role bindings for a subject.
type BindingReader interface {
	ListBindings(ctx context.Context, subject string) ([]Binding, error)
}
