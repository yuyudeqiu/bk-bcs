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
	"strings"
	"time"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/pkg/metrics"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/authorization"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

// AuthorizationStore reads local authorization bindings.
type AuthorizationStore struct{}

// NewAuthorizationStore creates a local authorization store.
func NewAuthorizationStore() *AuthorizationStore {
	return &AuthorizationStore{}
}

// ListBindings resolves all role bindings for a username.
func (s *AuthorizationStore) ListBindings(_ context.Context, subject string) ([]authorization.Binding, error) {
	bindings := make([]authorization.Binding, 0)
	err := GCoreDB.Table("bcs_user_resource_roles AS bindings").
		Select("bindings.resource_type, bindings.resource, roles.actions").
		Joins("JOIN bcs_roles AS roles ON roles.id = bindings.role_id").
		Where("bindings.subject = ?", subject).
		Scan(&bindings).Error
	return bindings, err
}

// EnsureDefaultRoles creates the built-in local roles when they do not exist.
func EnsureDefaultRoles() error {
	roles := []models.BcsRole{
		{Name: "manager", Actions: "*"},
		{Name: "viewer", Actions: strings.Join([]string{
			authorization.ActionHTTPGet,
			authorization.ActionProjectView,
			authorization.ActionClusterView,
			authorization.ActionClusterUse,
			authorization.ActionNamespaceView,
			authorization.ActionNamespaceList,
		}, ",")},
	}
	for i := range roles {
		if GetRole(roles[i].Name) != nil {
			continue
		}
		if err := CreateRole(&roles[i]); err != nil {
			return err
		}
	}
	return nil
}

// GetRole get bcsRole by roleName
func GetRole(roleName string) *models.BcsRole {
	start := time.Now()
	role := models.BcsRole{}
	GCoreDB.Where(&models.BcsRole{Name: roleName}).First(&role)
	if role.Name != "" {
		return &role
	}
	metrics.ReportMysqlSlowQueryMetrics("GetRole", metrics.Query, metrics.SucStatus, start)
	return nil
}

// CreateRole create role
func CreateRole(role *models.BcsRole) error {
	start := time.Now()
	err := GCoreDB.Create(role).Error
	if err != nil {
		metrics.ReportMysqlSlowQueryMetrics("CreateRole", metrics.Create, metrics.ErrStatus, start)
		return err
	}
	metrics.ReportMysqlSlowQueryMetrics("CreateRole", metrics.Create, metrics.SucStatus, start)
	return nil
}

// GetUrrByCondition Query BcsUserResourceRole by condition
func GetUrrByCondition(cond *models.BcsUserResourceRole) *models.BcsUserResourceRole {
	start := time.Now()
	urr := models.BcsUserResourceRole{}
	GCoreDB.Where(cond).First(&urr)
	if urr.ID != 0 {
		return &urr
	}
	metrics.ReportMysqlSlowQueryMetrics("GetUrrByCondition", metrics.Query, metrics.SucStatus, start)
	return nil
}

// CreateUserResourceRole create user-resource-role
func CreateUserResourceRole(urr *models.BcsUserResourceRole) error {
	start := time.Now()
	err := GCoreDB.Create(urr).Error
	if err != nil {
		metrics.ReportMysqlSlowQueryMetrics("CreateUserResourceRole", metrics.Create, metrics.ErrStatus, start)
		return err
	}
	metrics.ReportMysqlSlowQueryMetrics("CreateUserResourceRole", metrics.Create, metrics.SucStatus, start)
	return nil
}

// DeleteUserResourceRole delete user-resource-role
func DeleteUserResourceRole(urr *models.BcsUserResourceRole) error {
	start := time.Now()
	err := GCoreDB.Delete(urr).Error
	if err != nil {
		metrics.ReportMysqlSlowQueryMetrics("DeleteUserResourceRole", metrics.Delete, metrics.ErrStatus, start)
		return err
	}
	metrics.ReportMysqlSlowQueryMetrics("DeleteUserResourceRole", metrics.Delete, metrics.SucStatus, start)
	return nil
}
