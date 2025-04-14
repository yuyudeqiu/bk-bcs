package wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/common/headerkey"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/logging"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/errorx"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/util/tenant"
	middleauth "github.com/Tencent/bk-bcs/bcs-services/pkg/bcs-auth/middleware"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

var TenantClientWhiteList = map[string][]string{}

var NoCheckTenantMethod = []string{
	"BCSProject.CreateProject",
	"BCSProject.ListProjects",
	"BCSProject.ListAuthorizedProjects",
	"BCSProject.ListProjectsForIAM",
	"Healthz.Healthz",
	"Healthz.Ping",
	"Business.ListBusiness",
}

// SkipMethod skip method tenant validation
func SkipMethod(ctx context.Context, req server.Request) bool {
	for _, v := range NoCheckTenantMethod {
		if v == req.Method() {
			return true
		}
	}
	return false
}

func SkipTenantValidation(ctx context.Context, req server.Request, client string) bool {
	if len(client) == 0 {
		return false
	}
	for _, v := range TenantClientWhiteList[client] {
		if strings.HasPrefix(v, "*") || v == req.Method() {
			return true
		}
	}
	return false
}

// CheckUserResourceTenantAttrFunc xxx
func CheckUserResourceTenantAttrFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		if !tenant.IsMultiTenantEnabled() {
			return fn(ctx, req, rsp)
		}

		var (
			tenantId       = ""
			headerTenantId = GetHeaderTenantIdFromCtx(ctx)
			user, _        = middleauth.GetUserFromContext(ctx)
		)
		// exempt inner user
		if user.IsInner() {
			return fn(ctx, req, rsp)
		}
		// skip method tenant validation
		if SkipMethod(ctx, req) {
			return fn(ctx, req, rsp)
		}
		// exempt client
		if SkipTenantValidation(ctx, req, user.GetUsername()) {
			return fn(ctx, req, rsp)
		}

		// get tenant id
		if headerTenantId == "" {
			tenantId = user.GetTanantId()
		} else {
			if user.GetTanantId() != headerTenantId {
				tenantId = user.GetTanantId()
			} else {
				tenantId = headerTenantId
			}
		}

		// get resource tenant id
		resourceTenantId, err := getResourceTenantId(ctx, req)
		if err != nil {
			logging.Error("CheckUserResourceTenantAttrFunc getResourceTenantId failed, err: %s", err.Error())
			return err
		}
		logging.Info("CheckUserResourceTenantAttrFunc headerTenantId[%s] userTenantId[%s] tenantId[%s] resourceTenantId[%s]",
			headerTenantId, user.GetTanantId(), tenantId, resourceTenantId)
		if tenantId != resourceTenantId {
			return fmt.Errorf("user[%s] tenant[%s] not match resource tenant[%s]",
				user.GetUsername(), tenantId, resourceTenantId)
		}

		// 注入租户信息
		ctx = context.WithValue(ctx, headerkey.TenantIdKey, tenantId)

		return fn(ctx, req, rsp)
	}
}

// getResourceTenantId get tenant id
func getResourceTenantId(ctx context.Context, req server.Request) (string, error) {
	b, err := json.Marshal(req.Body())
	if err != nil {
		return "", err
	}
	r := &resourceID{}
	if err := json.Unmarshal(b, r); err != nil {
		return "", err
	}

	// 选择第一个非空值作为查询参数
	queryParam := r.ProjectIDOrCode
	if queryParam == "" {
		queryParam = r.ProjectCode
	}
	if queryParam == "" {
		queryParam = r.ProjectID
	}

	if queryParam == "" {
		return "", errorx.NewReadableErr(errorx.PermDeniedErr, "resource is empty")
	}

	p, err := store.GetModel().GetProject(ctx, queryParam)
	if err != nil {
		return "", errorx.NewReadableErr(errorx.ProjectNotExistsErr, "project is not exists")
	}

	r.ProjectID, r.ProjectCode = p.ProjectID, p.ProjectCode
	return p.TenantID, nil
}

// GetHeaderTenantIdFromCtx get tenantID from header
func GetHeaderTenantIdFromCtx(ctx context.Context) string {
	md, _ := metadata.FromContext(ctx)
	tenantId, _ := md.Get(headerkey.TenantIdKey)
	return tenantId
}
