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

func TestNew(t *testing.T) {
	reader := &stubBindingReader{}
	tests := []struct {
		name      string
		mode      string
		reader    BindingReader
		wantType  interface{}
		wantError bool
	}{
		{name: "empty defaults to none", mode: "", wantType: NoneAuthorizer{}},
		{name: "none is normalized", mode: " NONE ", wantType: NoneAuthorizer{}},
		{name: "local is normalized", mode: " LOCAL ", reader: reader, wantType: &LocalAuthorizer{}},
		{name: "local requires reader", mode: ModeLocal, wantError: true},
		{name: "unknown mode", mode: "iam", reader: reader, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authorizer, err := New(tt.mode, tt.reader)
			if tt.wantError {
				if err == nil {
					t.Fatal("New() returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatalf("New() returned error: %v", err)
			}
			switch tt.wantType.(type) {
			case NoneAuthorizer:
				if _, ok := authorizer.(NoneAuthorizer); !ok {
					t.Fatalf("New() returned %T, want NoneAuthorizer", authorizer)
				}
			case *LocalAuthorizer:
				if _, ok := authorizer.(*LocalAuthorizer); !ok {
					t.Fatalf("New() returned %T, want *LocalAuthorizer", authorizer)
				}
			}
		})
	}
}
