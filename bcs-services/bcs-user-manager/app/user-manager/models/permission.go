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

// BcsRole table
type BcsRole struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"unique;not null"`
	Actions   string `gorm:"not null" json:"actions"`
	CreatedAt time.Time
}

// BcsUserResourceRole table
type BcsUserResourceRole struct {
	ID           uint   `gorm:"primary_key"`
	Subject      string `gorm:"size:128;not null;unique_index:uk_subject_resource_role"`
	ResourceType string `gorm:"size:64;not null;unique_index:uk_subject_resource_role"`
	Resource     string `gorm:"not null;unique_index:uk_subject_resource_role"`
	RoleId       uint   `gorm:"not null;unique_index:uk_subject_resource_role"`
	CreatedAt    time.Time
}
