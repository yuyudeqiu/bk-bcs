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
	"fmt"
	"strings"
)

// New creates an Authorizer for the configured mode.
func New(mode string, bindings BindingReader) (Authorizer, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", ModeNone:
		return NoneAuthorizer{}, nil
	case ModeLocal:
		if bindings == nil {
			return nil, fmt.Errorf("authorization mode %q requires a binding reader", mode)
		}
		return NewLocalAuthorizer(bindings), nil
	default:
		return nil, fmt.Errorf("unsupported authorization mode %q", mode)
	}
}
