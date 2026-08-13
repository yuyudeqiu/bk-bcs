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

package usermanager

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/Tencent/bk-bcs/bcs-common/common"
	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	bcshttp "github.com/Tencent/bk-bcs/bcs-common/common/http"
	"github.com/Tencent/bk-bcs/bcs-common/common/http/httpserver"
	"github.com/Tencent/bk-bcs/bcs-common/common/ssl"
	"github.com/Tencent/bk-bcs/bcs-common/pkg/i18n"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"

	i18n2 "github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/pkg/i18n"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/authorization"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/job/activity"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/cache"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/v1http"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/v1http/permission"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/v3http"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/utils"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/config"
)

var (
	// ErrHTTPServerNotInit server not init
	ErrHTTPServerNotInit = errors.New("UserManager server not init")
)

// UserManager http interface of user-manager
type UserManager struct {
	config   *config.UserMgrConfig
	httpServ *httpserver.HttpServer

	EtcdRegistry registry.Registry

	authorizer  authorization.Authorizer
	permService *permission.PermVerifyClient
}

// NewUserManager creates an UserManager object
func NewUserManager(conf *config.UserMgrConfig) *UserManager {
	userManager := &UserManager{
		config:   conf,
		httpServ: httpserver.NewIPv6HttpServer(conf.Port, conf.Address, conf.IPv6Address, conf.Sock),
	}

	if conf.ServCert.IsSSL {
		userManager.httpServ.SetSsl(conf.ServCert.CAFile, conf.ServCert.CertFile, conf.ServCert.KeyFile,
			conf.ServCert.CertPasswd)
	}

	userManager.httpServ.SetInsecureServer(conf.InsecureAddress, conf.InsecurePort)

	return userManager
}

// Start entry point for user-manager
func (u *UserManager) Start() error {
	// init redis
	if err := cache.InitRedis(u.config); err != nil {
		return err
	}

	if err := SetupStore(u.config); err != nil {
		return err
	}

	// 定时清理操作记录
	go func() {
		err := activity.IntervalDeleteActivity(context.Background())
		if err != nil {
			blog.Errorf("IntervalDeleteActivity failed: %v", err)
		}
	}()

	err := u.initUserManagerServer()
	if err != nil {
		blog.Errorf("initUserManagerServer failed: %v", err)
		return err
	}

	// usermanager api
	v1http.InitV1Routers(u.httpServ.NewWebService("/usermanager", nil), u.permService, u.authorizer)
	v3http.InitV3Routers(u.httpServ.NewWebService("/usermanager/v3", nil), u.authorizer)

	router := u.httpServ.GetRouter()
	webContainer := u.httpServ.GetWebContainer()
	// set recover handler
	webContainer.RecoverHandler(responseOnRecover)
	webContainer.DoNotRecover(false)

	// handle user and cluster manager request
	router.Handle("/usermanager/{sub_path:.*}", webContainer)

	if err := u.httpServ.ListenAndServeMux(u.config.VerifyClientTLS); err != nil {
		return fmt.Errorf("http ListenAndServe error %s", err.Error())
	}

	return nil
}

// Filter authenticate the request
func Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// first authenticate the request, only admin user be allowed
	auth := utils.Authenticate(req.Request)
	if !auth {
		_ = resp.WriteHeaderAndEntity(http.StatusUnauthorized, bcshttp.APIRespone{
			Result:  false,
			Code:    common.BcsErrApiUnauthorized,
			Message: "must provide admin token to request with websocket",
			Data:    nil,
		})
		return
	}

	chain.ProcessFilter(req, resp)
}

func (u *UserManager) initPermService() error {
	u.permService = permission.NewPermVerifyClient(u.authorizer)

	return nil
}

func (u *UserManager) initAuthorizer() error {
	authorizer, err := authorization.New(u.config.Authorization.Mode, sqlstore.NewAuthorizationStore())
	if err != nil {
		return err
	}
	u.authorizer = authorizer
	return nil
}

func (u *UserManager) initEtcdRegistry() error {
	if !u.config.EtcdConfig.Feature {
		return fmt.Errorf("etcd feature is off")
	}

	if len(u.config.EtcdConfig.Address) == 0 {
		errMsg := fmt.Errorf("etcdServers invalid")
		return errMsg
	}
	servers := strings.Split(u.config.EtcdConfig.Address, ";")

	var (
		secureEtcd bool
		etcdTLS    *tls.Config
		err        error
	)

	if len(u.config.EtcdConfig.CA) != 0 && len(u.config.EtcdConfig.Cert) != 0 && len(u.config.EtcdConfig.Key) != 0 {
		secureEtcd = true

		etcdTLS, err = ssl.ClientTslConfVerity(u.config.EtcdConfig.CA, u.config.EtcdConfig.Cert,
			u.config.EtcdConfig.Key, "")
		if err != nil {
			return err
		}
	}

	u.EtcdRegistry = etcd.NewRegistry(
		registry.Addrs(servers...),
		registry.Secure(secureEtcd),
		registry.TLSConfig(etcdTLS),
	)
	if err := u.EtcdRegistry.Init(); err != nil {
		return err
	}

	return nil
}

// initI18n init i18n
func (u *UserManager) initI18n() {
	i18n.Instance()
	// 加载翻译文件路径
	i18n.SetPath([]embed.FS{i18n2.Assets})
	// 设置默认语言
	// 默认是 zh
	i18n.SetLanguage("zh")
}

func (u *UserManager) initUserManagerServer() error {
	var err error
	err = u.initEtcdRegistry()
	if err != nil {
		return err
	}

	err = u.initAuthorizer()
	if err != nil {
		return err
	}

	err = u.initPermService()
	if err != nil {
		return err
	}

	u.initI18n()

	return nil
}

// responseOnRecover response on recover
func responseOnRecover(panicReason interface{}, httpWriter http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("recover from panic situation: - %v\r\n", panicReason))
	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	blog.Error(buffer.String())
	httpWriter.WriteHeader(http.StatusInternalServerError)
	httpWriter.Write([]byte(`{"code": 500, "message": "server error"}`))
}
