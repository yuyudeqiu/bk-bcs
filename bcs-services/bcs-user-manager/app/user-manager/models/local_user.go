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

package models

import "time"

const (
	// LocalUserStatusActive means the local user is allowed to log in.
	LocalUserStatusActive = "active"
	// LocalUserStatusDisabled means the local user is not allowed to log in.
	LocalUserStatusDisabled = "disabled"
)

// LocalUser is an account authenticated directly by bcs-user-manager.
// Subject is stable and must be used for authorization bindings instead of Username.
type LocalUser struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	Subject      string     `json:"subject" gorm:"type:varchar(128);not null;unique_index"`
	Username     string     `json:"username" gorm:"type:varchar(64);not null;unique_index"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null"`
	Status       string     `json:"status" gorm:"type:varchar(16);not null"`
	IsAdmin      bool       `json:"is_admin" gorm:"not null;default:false"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp null;default:null"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:timestamp null;default:null"`
	DeletedAt    *time.Time `json:"-" gorm:"type:timestamp null;default:null"`
}
