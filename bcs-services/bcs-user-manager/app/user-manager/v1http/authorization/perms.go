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
	"github.com/Tencent/bk-bcs/bcs-common/common"
	restful "github.com/emicklei/go-restful/v3"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/pkg/constant"
	coreauth "github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/authorization"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/v1http/auth"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/utils"
)

// PermRequest perm request
type PermRequest struct {
	ActionIDs []string      `json:"action_ids"`
	PermCtx   *auth.PermCtx `json:"perm_ctx"`
}

// GetPerms get perm
func GetPerms(authorizer coreauth.Authorizer) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		form := PermRequest{}
		_ = request.ReadEntity(&form)
		err := utils.Validate.Struct(&form)
		if err != nil {
			_ = response.WriteHeaderAndEntity(400, utils.FormatValidationError(err))
			return
		}

		user := getCurrentUser(request)
		if user == nil {
			utils.WriteUnauthorizedError(response, common.BcsErrApiUnauthorized, "user is not valid")
			return
		}

		result := make(map[string]bool, len(form.ActionIDs))
		for _, actionID := range form.ActionIDs {
			decision, authErr := authorizer.Authorize(request.Request.Context(), coreauth.Request{
				Subject:   user.Name,
				Superuser: user.IsAdmin(),
				Action:    actionID,
				Resource:  auth.ResourceFromPermCtx(form.PermCtx),
			})
			if authErr != nil {
				utils.WriteServerError(response, common.BcsErrApiBadRequest, authErr.Error())
				return
			}
			result[actionID] = decision.Allowed
		}

		data := utils.CreateResponseData(nil, "success", map[string]interface{}{"perms": result})
		_, _ = response.Write([]byte(data))
	}
}

// GetPermByActionID get perm by action id
func GetPermByActionID(authorizer coreauth.Authorizer) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		actionID := request.PathParameter("action_id")
		form := PermRequest{}
		_ = request.ReadEntity(&form)
		err := utils.Validate.Struct(&form)
		if err != nil {
			_ = response.WriteHeaderAndEntity(400, utils.FormatValidationError(err))
			return
		}
		if form.PermCtx != nil && form.PermCtx.ResourceType == "" {
			form.PermCtx.ResourceType = auth.GetResourceTypeFromAction(actionID)
		}

		user := getCurrentUser(request)
		if user == nil {
			utils.WriteUnauthorizedError(response, common.BcsErrApiUnauthorized, "user is not valid")
			return
		}

		decision, err := authorizer.Authorize(request.Request.Context(), coreauth.Request{
			Subject:   user.Name,
			Superuser: user.IsAdmin(),
			Action:    actionID,
			Resource:  auth.ResourceFromPermCtx(form.PermCtx),
		})
		if err != nil {
			utils.WriteServerError(response, common.BcsErrApiBadRequest, err.Error())
			return
		}

		data := utils.CreateResponseData(nil, "success", map[string]interface{}{
			"perms": map[string]interface{}{actionID: decision.Allowed, "apply_url": ""}})
		_, _ = response.Write([]byte(data))
	}
}

func getCurrentUser(request *restful.Request) *models.BcsUser {
	user, _ := request.Attribute(constant.CurrentUserAttr).(*models.BcsUser)
	return user
}
