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

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

type fakeLocalUserStore struct {
	users     map[string]*models.LocalUser
	getErr    error
	createErr error
	creates   int
}

func newFakeLocalUserStore() *fakeLocalUserStore {
	return &fakeLocalUserStore{users: make(map[string]*models.LocalUser)}
}

func (s *fakeLocalUserStore) GetByUsername(_ context.Context, username string) (*models.LocalUser, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.users[username], nil
}

func (s *fakeLocalUserStore) Create(_ context.Context, user *models.LocalUser) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.creates++
	s.users[user.Username] = user
	return nil
}

func TestEnsureBootstrapAdminDisabled(t *testing.T) {
	store := newFakeLocalUserStore()
	user, created, err := EnsureBootstrapAdmin(context.Background(), store, "", "")
	if err != nil || user != nil || created {
		t.Fatalf("EnsureBootstrapAdmin() = (%v, %v, %v), want disabled", user, created, err)
	}
}

func TestEnsureBootstrapAdminCreatesOnce(t *testing.T) {
	store := newFakeLocalUserStore()
	ctx := context.Background()

	createdUser, created, err := EnsureBootstrapAdmin(ctx, store, " admin ", "local-password")
	if err != nil {
		t.Fatalf("EnsureBootstrapAdmin() error = %v", err)
	}
	if !created || store.creates != 1 {
		t.Fatalf("EnsureBootstrapAdmin() created = %v, creates = %d", created, store.creates)
	}
	if createdUser.Username != "admin" || !createdUser.IsAdmin || createdUser.Status != models.LocalUserStatusActive {
		t.Fatalf("EnsureBootstrapAdmin() user = %#v", createdUser)
	}
	if !strings.HasPrefix(createdUser.Subject, "user:") {
		t.Fatalf("EnsureBootstrapAdmin() subject = %q", createdUser.Subject)
	}
	if !VerifyPassword(createdUser.PasswordHash, "local-password") {
		t.Fatal("EnsureBootstrapAdmin() password hash does not verify")
	}

	existingUser, created, err := EnsureBootstrapAdmin(ctx, store, "admin", "different-password")
	if err != nil {
		t.Fatalf("second EnsureBootstrapAdmin() error = %v", err)
	}
	if created || store.creates != 1 || existingUser != createdUser {
		t.Fatalf("second EnsureBootstrapAdmin() created = %v, creates = %d", created, store.creates)
	}
	if VerifyPassword(existingUser.PasswordHash, "different-password") {
		t.Fatal("second EnsureBootstrapAdmin() unexpectedly reset the password")
	}
}

func TestEnsureBootstrapAdminRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{name: "username only", username: "admin"},
		{name: "password only", password: "local-password"},
		{name: "short password", username: "admin", password: "short"},
		{name: "long username", username: strings.Repeat("a", maxUsernameLength+1), password: "local-password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := EnsureBootstrapAdmin(context.Background(), newFakeLocalUserStore(), tt.username, tt.password)
			if err == nil {
				t.Fatal("EnsureBootstrapAdmin() expected an error")
			}
		})
	}
}

func TestEnsureBootstrapAdminPropagatesStoreErrors(t *testing.T) {
	store := newFakeLocalUserStore()
	store.getErr = errors.New("get failed")
	if _, _, err := EnsureBootstrapAdmin(context.Background(), store, "admin", "local-password"); err == nil {
		t.Fatal("EnsureBootstrapAdmin() did not propagate get error")
	}

	store.getErr = nil
	store.createErr = errors.New("create failed")
	if _, _, err := EnsureBootstrapAdmin(context.Background(), store, "admin", "local-password"); err == nil {
		t.Fatal("EnsureBootstrapAdmin() did not propagate create error")
	}
}
