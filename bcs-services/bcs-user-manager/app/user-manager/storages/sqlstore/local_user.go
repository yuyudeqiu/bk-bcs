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

package sqlstore

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

// LocalUserStore stores accounts authenticated by bcs-user-manager.
type LocalUserStore struct {
	db *gorm.DB
}

// NewLocalUserStore creates a local user store.
func NewLocalUserStore(db *gorm.DB) *LocalUserStore {
	return &LocalUserStore{db: db}
}

// GetByUsername returns nil when the account does not exist.
func (s *LocalUserStore) GetByUsername(_ context.Context, username string) (*models.LocalUser, error) {
	var user models.LocalUser
	err := s.db.Where("username = ?", username).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create persists a local user.
func (s *LocalUserStore) Create(_ context.Context, user *models.LocalUser) error {
	return s.db.Create(user).Error
}
