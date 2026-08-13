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

package authentication

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("local-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "local-password" {
		t.Fatal("HashPassword() stored the plaintext password")
	}
	if !VerifyPassword(hash, "local-password") {
		t.Fatal("VerifyPassword() rejected the correct password")
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Fatal("VerifyPassword() accepted an incorrect password")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "minimum", password: "12345678"},
		{name: "too short", password: "1234567", wantErr: true},
		{name: "maximum", password: string(make([]byte, maxPasswordLength))},
		{name: "too long", password: string(make([]byte, maxPasswordLength+1)), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
