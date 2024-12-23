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

package RegisterDiscover

import (
	"time"
)

// DiscoverEvent if servers chenged, will create a discover event
type DiscoverEvent struct { //
	Err    error
	Key    string
	Server []string
	Nodes  []string
}

// RegDiscover is data struct of register-discover
type RegDiscover struct {
	rdServer RegDiscvServer
}

// NewRegDiscover used to create a object of RegDiscover
func NewRegDiscover(serv string) *RegDiscover {
	regDiscv := &RegDiscover{
		rdServer: nil,
	}

	regDiscv.rdServer = RegDiscvServer(NewZkRegDiscv(serv, time.Second*60))

	return regDiscv
}

// NewRegDiscoverEx used to create a object of RegDiscover
func NewRegDiscoverEx(serv string, timeOut time.Duration) *RegDiscover {
	regDiscv := &RegDiscover{
		rdServer: nil,
	}

	regDiscv.rdServer = RegDiscvServer(NewZkRegDiscv(serv, timeOut))

	return regDiscv
}

// Start the register and discover service
func (rd *RegDiscover) Start() error {
	return rd.rdServer.Start()
}

// Stop the register and discover service
func (rd *RegDiscover) Stop() error {
	return rd.rdServer.Stop()
}

// RegisterService used to write service info into register and discover server
// key is the index of registered service
// data is the service information
func (rd *RegDiscover) RegisterService(key string, data []byte) error {
	return rd.rdServer.Register(key, data)
}

// RegisterAndWatchService register service info into register-discover platform
// and then watch the service info, if not exist, then register again
// key is the index of registered service
// data is the service information
func (rd *RegDiscover) RegisterAndWatchService(key string, data []byte) error {
	return rd.rdServer.RegisterAndWatch(key, data)
}

// DiscoverService used to discover the service that registered in `key`
func (rd *RegDiscover) DiscoverService(key string) (<-chan *DiscoverEvent, error) {
	return rd.rdServer.Discover(key)
}

// DiscoverNodes discover nodes from path
func (rd *RegDiscover) DiscoverNodes(path string) (*DiscoverEvent, error) {
	return rd.rdServer.DiscoverNodes(path)
}

// DiscoverNodesV2 discover nodes from path V2
func (rd *RegDiscover) DiscoverNodesV2(path string) (*DiscoverEvent, error) {
	return rd.rdServer.DiscoverNodesV2(path)
}
