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
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

const maxUsernameLength = 64

// LocalUserStore persists local user accounts.
type LocalUserStore interface {
	GetByUsername(ctx context.Context, username string) (*models.LocalUser, error)
	Create(ctx context.Context, user *models.LocalUser) error
}

// ValidatePassword validates passwords before they are hashed.
func ValidatePassword(password string) error {
	length := len([]byte(password))
	if length < minPasswordLength {
		return fmt.Errorf("password must contain at least %d bytes", minPasswordLength)
	}
	if length > maxPasswordLength {
		return fmt.Errorf("password must contain at most %d bytes", maxPasswordLength)
	}
	return nil
}

// EnsureBootstrapAdmin creates the initial local administrator once.
// Empty username and password disable bootstrap. An existing account is never modified.
func EnsureBootstrapAdmin(ctx context.Context, store LocalUserStore, username, password string) (
	*models.LocalUser, bool, error) {
	username = strings.TrimSpace(username)
	if username == "" && password == "" {
		return nil, false, nil
	}
	if username == "" || password == "" {
		return nil, false, fmt.Errorf("bootstrap admin username and password must be configured together")
	}
	if utf8.RuneCountInString(username) > maxUsernameLength {
		return nil, false, fmt.Errorf("username must contain at most %d characters", maxUsernameLength)
	}
	if err := ValidatePassword(password); err != nil {
		return nil, false, err
	}

	existing, err := store.GetByUsername(ctx, username)
	if err != nil {
		return nil, false, fmt.Errorf("get bootstrap admin: %w", err)
	}
	if existing != nil {
		return existing, false, nil
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, false, fmt.Errorf("hash bootstrap admin password: %w", err)
	}
	user := &models.LocalUser{
		Subject:      "user:" + uuid.NewString(),
		Username:     username,
		PasswordHash: passwordHash,
		Status:       models.LocalUserStatusActive,
		IsAdmin:      true,
	}
	if err := store.Create(ctx, user); err != nil {
		return nil, false, fmt.Errorf("create bootstrap admin: %w", err)
	}
	return user, true, nil
}
