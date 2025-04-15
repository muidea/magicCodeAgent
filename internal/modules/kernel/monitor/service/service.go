package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/monitor/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/monitor/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type Monitor struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Monitor
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Monitor) *Monitor {
	ptr := &Monitor{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Monitor) RegisterRoute() {
	s.routeRegistry.SetApiVersion(tc.ApiVersion)
	s.routeRegistry.AddHandler(common.FilterWorkspace, engine.GET, s.filterWorkspace)
	s.routeRegistry.AddHandler(common.AddWorkspace, engine.POST, s.addWorkspace)
	s.routeRegistry.AddHandler(common.RemoveWorkspace, engine.DELETE, s.delWorkspace)
	s.routeRegistry.AddHandler(common.EnableWorkspace, engine.PUT, s.enableWorkspace)
	s.routeRegistry.AddHandler(common.DisableWorkspace, engine.PUT, s.disableWorkspace)
	s.routeRegistry.AddHandler(common.RefreshWorkspace, engine.GET, s.refreshWorkspace)
	s.routeRegistry.AddHandler(common.FilterWorkspaceFile, engine.GET, s.filterFile)
}

func (s *Monitor) filterWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
}

func (s *Monitor) addWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceResult{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}

	result.Value, result.Error = s.bizPtr.AddWorkspace(param)
}

func (s *Monitor) delWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &cd.Result{}
	defer fn.PackageHTTPResponse(res, result)
	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.RemoveWorkspace(workspaceName)
}

func (s *Monitor) enableWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &cd.Result{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	workspaceName := param.Name
	result.Error = s.bizPtr.EnableWorkspace(workspaceName)
}

func (s *Monitor) disableWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &cd.Result{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	workspaceName := param.Name
	result.Error = s.bizPtr.DisableWorkspace(workspaceName)
}

func (s *Monitor) refreshWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &cd.Result{}
	defer fn.PackageHTTPResponse(res, result)
	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.RefreshWorkspace(workspaceName)
}

/*
| 请求参数**（Query Parameters）
| `name`       | string | 否   | 文件名模糊匹配（支持通配符 `*`，如 `*report*` 匹配包含 "report" 的文件） |
| `path`       | string | 否   | 文件路径前缀匹配（如 `/docs/2025` 匹配该路径下的所有文件）              |
| `extension`  | string | 否   | 文件扩展名精确匹配（如 `pdf`，不区分大小写）                          |
*/
func (s *Monitor) filterFile(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.FileListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	filter := &common.FileFilter{}
	filter.Parse(req)
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileList(filter.Workspace, filter)
}
