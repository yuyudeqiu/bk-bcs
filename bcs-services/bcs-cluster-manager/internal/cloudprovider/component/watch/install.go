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

// Package watch xxx
package watch

import (
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/component"
	cmoptions "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/options"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/remote/install"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/remote/install/helm"
)

// GetWatchInstaller watch installer
func GetWatchInstaller(projectID string, namespace string) (install.Installer, error) {
	op := cmoptions.GetGlobalCMOptions()

	return component.GetComponentInstaller(component.InstallOptions{
		InstallType:      helm.Helm.String(),
		ProjectID:        projectID,
		ChartName:        op.ComponentDeploy.Watch.ChartName,
		ReleaseNamespace: namespace,
		ReleaseName:      op.ComponentDeploy.Watch.ReleaseName,
		IsPublicRepo:     op.ComponentDeploy.Watch.IsPublicRepo,
	})
}
