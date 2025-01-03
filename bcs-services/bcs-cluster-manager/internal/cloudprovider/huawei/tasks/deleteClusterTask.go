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

package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/huawei/api"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/huawei/business"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/cloudprovider/utils"
	icommon "github.com/Tencent/bk-bcs/bcs-services/bcs-cluster-manager/internal/common"
)

// DeleteClusterTask delete cluster task
func DeleteClusterTask(taskID string, stepName string) error { // nolint
	start := time.Now()
	// get task and task current step
	state, step, err := cloudprovider.GetTaskStateAndCurrentStep(taskID, stepName)
	if err != nil {
		return err
	}
	// previous step successful when retry task
	if step == nil {
		blog.Infof("DeleteClusterTask[%s]: current step[%s] successful and skip", taskID, stepName)
		return nil
	}
	blog.Infof("DeleteClusterTask[%s]: task %s run step %s, system: %s, old state: %s, params %v",
		taskID, taskID, stepName, step.System, step.Status, step.Params)

	// step login started here
	clusterID := step.Params[cloudprovider.ClusterIDKey.String()]
	cloudID := step.Params[cloudprovider.CloudIDKey.String()]
	deleteMode := step.Params[cloudprovider.DeleteModeKey.String()]
	clusterStatus := step.Params[cloudprovider.LastClusterStatus.String()]

	// only support retain mode
	if deleteMode != cloudprovider.Retain.String() {
		deleteMode = cloudprovider.Retain.String()
	}

	dependInfo, err := cloudprovider.GetClusterDependBasicInfo(cloudprovider.GetBasicInfoReq{
		ClusterID: clusterID,
		CloudID:   cloudID,
	})
	if err != nil {
		blog.Errorf("DeleteClusterTask[%s]: GetClusterDependBasicInfo for cluster %s "+
			"in task %s step %s failed, %s", taskID, clusterID, taskID, stepName, err.Error())
		retErr := fmt.Errorf("get cloud/project information failed, %s", err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)
		return retErr
	}

	client, err := api.NewCceClient(dependInfo.CmOption)
	if err != nil {
		return err
	}

	ctx := cloudprovider.WithTaskIDForContext(context.Background(), taskID)

	blog.Infof("DeleteTKEClusterTask[%s]  clusterInfo: %v", taskID, dependInfo.Cluster.GetStatus())

	if (clusterStatus == icommon.StatusCreateClusterFailed ||
		clusterStatus == icommon.StatusDeleteClusterFailed) && dependInfo.Cluster.GetSystemID() != "" {
		nodes, listErr := client.ListClusterNodes(dependInfo.Cluster.SystemID)
		if listErr != nil {
			return listErr
		}

		ids, delErr := business.DeleteClusterInstance(client, dependInfo.Cluster.SystemID, nodes)
		if delErr != nil {
			blog.Errorf("DeleteClusterTask[%s] DeleteClusterInstance failed: %v", taskID, delErr)
		} else {
			blog.Infof("DeleteClusterTask[%s] DeleteClusterInstance success: %v", taskID, ids)
		}

		nodeIds := make([]string, 0)
		for _, node := range nodes {
			nodeIds = append(nodeIds, *node.Metadata.Uid)
		}

		err = business.CheckClusterDeletedNodes(ctx, dependInfo, nodeIds)
		if err != nil {
			blog.Errorf("DeleteClusterTask[%s] CheckClusterDeletedNodes failed: %v", taskID, err)
		}
	}

	err = business.DeleteTkeClusterByClusterId(ctx, dependInfo.CmOption, dependInfo.Cluster.SystemID, deleteMode)
	if err != nil {
		blog.Errorf("DeleteClusterTask[%s]: task[%s] step[%s] call huawei DeleteCluster failed: %v",
			taskID, taskID, stepName, err)
		retErr := fmt.Errorf("call huawei DeleteTKECluster failed: %s", err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)

		_ = cloudprovider.UpdateClusterErrMessage(clusterID, fmt.Sprintf("delete cluster[%s] failed: %v",
			dependInfo.Cluster.GetClusterID(), err.Error()))

		return retErr
	}

	blog.Infof("DeleteClusterTask[%s]: task %s DeleteCluster[%s] successful",
		taskID, taskID, dependInfo.Cluster.SystemID)

	if err := state.UpdateStepSucc(start, stepName); err != nil {
		blog.Errorf("DeleteClusterTask[%s]: task %s %s update to storage fatal", taskID, taskID, stepName)
		return err
	}

	return nil
}

// CleanClusterDBInfoTask clean cluster DB info
func CleanClusterDBInfoTask(taskID string, stepName string) error {
	start := time.Now()
	// get task and task current step
	state, step, err := cloudprovider.GetTaskStateAndCurrentStep(taskID, stepName)
	if err != nil {
		return err
	}
	// previous step successful when retry task
	if step == nil {
		blog.Infof("CleanClusterDBInfoTask[%s]: current step[%s] successful and skip", taskID, stepName)
		return nil
	}
	blog.Infof("CleanClusterDBInfoTask[%s]: task %s run step %s, system: %s, old state: %s, params %v",
		taskID, taskID, stepName, step.System, step.Status, step.Params)

	// step login started here
	clusterID := step.Params[cloudprovider.ClusterIDKey.String()]
	cluster, err := cloudprovider.GetStorageModel().GetCluster(context.Background(), clusterID)
	if err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: get cluster for %s failed", taskID, clusterID)
		retErr := fmt.Errorf("get cluster information failed, %s", err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)
		return retErr
	}

	// delete cluster autoscalingOption
	err = cloudprovider.GetStorageModel().DeleteAutoScalingOption(context.Background(), cluster.ClusterID)
	if err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: clean cluster[%s] "+
			"autoscalingOption failed: %v", taskID, cluster.ClusterID, err)
	}

	// delete nodes
	err = cloudprovider.GetStorageModel().DeleteNodesByClusterID(context.Background(), cluster.ClusterID)
	if err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: delete nodes for %s failed", taskID, clusterID)
		retErr := fmt.Errorf("delete node for %s failed, %s", clusterID, err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)
		return retErr
	}
	blog.Infof("CleanClusterDBInfoTask[%s]: delete nodes for cluster[%s] in DB successful", taskID, clusterID)

	// delete nodeGroup
	err = cloudprovider.GetStorageModel().DeleteNodeGroupByClusterID(context.Background(), cluster.ClusterID)
	if err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: delete nodeGroups for %s failed", taskID, clusterID)
		retErr := fmt.Errorf("delete nodeGroups for %s failed, %s", clusterID, err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)
		return retErr
	}
	blog.Infof("CleanClusterDBInfoTask[%s]: delete nodeGroups for cluster[%s] in DB successful", taskID, clusterID)

	// delete cluster
	cluster.Status = icommon.StatusDeleting
	err = cloudprovider.GetStorageModel().UpdateCluster(context.Background(), cluster)
	if err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: delete cluster for %s failed", taskID, clusterID)
		retErr := fmt.Errorf("delete cluster for %s failed, %s", clusterID, err.Error())
		_ = state.UpdateStepFailure(start, stepName, retErr)
		return retErr
	}
	blog.Infof("CleanClusterDBInfoTask[%s]: delete cluster[%s] in DB successful", taskID, clusterID)

	utils.SyncDeletePassCCCluster(taskID, cluster)
	_ = utils.DeleteClusterCredentialInfo(cluster.ClusterID)

	if err := state.UpdateStepSucc(start, stepName); err != nil {
		blog.Errorf("CleanClusterDBInfoTask[%s]: task %s %s update to storage fatal", taskID, taskID, stepName)
		return err
	}
	return nil
}
