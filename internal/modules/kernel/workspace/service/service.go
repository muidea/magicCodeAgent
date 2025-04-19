package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type Workspace struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Workspace
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Workspace) *Workspace {
	ptr := &Workspace{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Workspace) RegisterRoute() {
	s.routeRegistry.SetApiVersion(tc.ApiVersion)
	s.routeRegistry.AddHandler(common.CreateWorkspace, engine.POST, s.createWorkspace)
	s.routeRegistry.AddHandler(common.FilterWorkspace, engine.GET, s.filterWorkspace)
	s.routeRegistry.AddHandler(common.AddWorkspace, engine.POST, s.addWorkspace)
	s.routeRegistry.AddHandler(common.UpdateWorkspace, engine.PUT, s.updateWorkspace)
	s.routeRegistry.AddHandler(common.RemoveWorkspace, engine.DELETE, s.delWorkspace)
	s.routeRegistry.AddHandler(common.EnableWorkspace, engine.PUT, s.enableWorkspace)
	s.routeRegistry.AddHandler(common.DisableWorkspace, engine.PUT, s.disableWorkspace)
	s.routeRegistry.AddHandler(common.RefreshWorkspace, engine.GET, s.refreshWorkspace)
	s.routeRegistry.AddHandler(common.FilterWorkspaceFile, engine.GET, s.filterFile)
	s.routeRegistry.AddHandler(common.FilterWorkspaceFileTree, engine.GET, s.filterFileTree)
}

func (s *Workspace) filterWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
}

func (s *Workspace) createWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}

	_, result.Error = s.bizPtr.CreateWorkspace(param)
	if result.Error == nil {
		// 创建成功后刷新工作区列表
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) addWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}

	_, result.Error = s.bizPtr.AddWorkspace(param)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) updateWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
	}
	result.Error = s.bizPtr.UpdateWorkspace(param)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) delWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)
	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.RemoveWorkspace(workspaceName)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) enableWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
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
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) disableWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
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
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) refreshWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)
	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.RefreshWorkspace(workspaceName)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

/*
| 请求参数**（Query Parameters）
| `name`       | string | 否   | 文件名模糊匹配（支持通配符 `*`，如 `*report*` 匹配包含 "report" 的文件） |
| `path`       | string | 否   | 文件路径前缀匹配（如 `/docs/2025` 匹配该路径下的所有文件）              |
| `extension`  | string | 否   | 文件扩展名精确匹配（如 `pdf`，不区分大小写）                          |
*/
func (s *Workspace) filterFile(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.FileListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	filter := &common.FileFilter{}
	filter.Parse(req)
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileList(filter.Workspace, filter)
}

func (s *Workspace) filterFileTree(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.FileTreeListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	workspace := req.URL.Query().Get("workspace")
	parentPath := req.URL.Query().Get("parentPath")
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileTree(workspace, parentPath)
}
