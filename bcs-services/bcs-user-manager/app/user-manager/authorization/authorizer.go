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
)

// Resource identifies the object affected by an action.
type Resource struct {
	Type       string
	ID         string
	Attributes map[string]string
}

// Request contains the information required to make an authorization decision.
type Request struct {
	Subject  string
	Action   string
	Resource Resource
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
