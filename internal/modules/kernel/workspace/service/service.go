package service

import (
	"context"
	"errors"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/biz"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type Workspace struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Workspace

	sseHolderRegistry *sse.HolderRegistry
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Workspace) *Workspace {
	ptr := &Workspace{
		routeRegistry:     routeRegistry,
		bizPtr:            bizPtr,
		sseHolderRegistry: sse.CreateHolderRegistry(),
	}
	return ptr
}

func (s *Workspace) RegisterRoute() {
	s.routeRegistry.SetApiVersion(tc.ApiVersion)
	s.routeRegistry.AddHandler(imkwc.QueryWorkspace, engine.GET, s.queryWorkspace)
	s.routeRegistry.AddHandler(imkwc.CreateWorkspace, engine.POST, s.createWorkspace)
	s.routeRegistry.AddHandler(imkwc.DestroyWorkspace, engine.DELETE, s.destroyWorkspace)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspace, engine.GET, s.filterWorkspace)
	s.routeRegistry.AddHandler(imkwc.AddWorkspace, engine.POST, s.addWorkspace)
	s.routeRegistry.AddHandler(imkwc.UpdateWorkspace, engine.PUT, s.updateWorkspace)
	s.routeRegistry.AddHandler(imkwc.RemoveWorkspace, engine.DELETE, s.delWorkspace)
	s.routeRegistry.AddHandler(imkwc.EnableWorkspace, engine.PUT, s.enableWorkspace)
	s.routeRegistry.AddHandler(imkwc.DisableWorkspace, engine.PUT, s.disableWorkspace)
	s.routeRegistry.AddHandler(imkwc.UpdateWorkspaceStage, engine.PUT, s.updateWorkspaceStage)
	s.routeRegistry.AddHandler(imkwc.RefreshWorkspace, engine.GET, s.refreshWorkspace)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceFile, engine.GET, s.filterFile)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceFileTree, engine.GET, s.filterFileTree)
	s.routeRegistry.AddHandler(imkwc.QueryWorkspaceFileContent, engine.GET, s.queryFileContent)
	s.routeRegistry.AddHandler(imkwc.ChatWorkspaceContentSSE, engine.GET, s.chatWorkspaceContentSSE)
	s.routeRegistry.AddHandler(imkwc.ChatWorkspaceContentMessage, engine.POST, s.chatWorkspaceContentMessage)
	s.routeRegistry.AddHandler(imkwc.ChatWorkspaceContentTerminate, engine.POST, s.chatWorkspaceContentTerminate)
	s.routeRegistry.AddHandler(imkwc.ChatWorkspaceContentHistory, engine.GET, s.chatWorkspaceContentHistory)
}

func (s *Workspace) filterWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
}

func (s *Workspace) queryWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceName := req.URL.Query().Get("name")
	result.Value, result.Error = s.bizPtr.QueryWorkspace(workspaceName)
}

func (s *Workspace) createWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
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

func (s *Workspace) destroyWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.DestroyWorkspace(workspaceName)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) addWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
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
	result := &imkwc.WorkspaceListResult{}
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
	result := &imkwc.WorkspaceListResult{}
	defer fn.PackageHTTPResponse(res, result)
	workspaceName := req.URL.Query().Get("name")
	result.Error = s.bizPtr.RemoveWorkspace(workspaceName)
	if result.Error == nil {
		result.Values, result.Error = s.bizPtr.FilterWorkspace(nil)
	}
}

func (s *Workspace) enableWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
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
	result := &imkwc.WorkspaceListResult{}
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

func (s *Workspace) updateWorkspaceStage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceResult{}
	defer fn.PackageHTTPResponse(res, result)
	param := &config.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	if !common.ValidStage(param.Stage) {
		result.Error = cd.NewError(cd.IllegalParam, "stage can not be empty")
		return
	}
	result.Error = s.bizPtr.UpdateWorkspace(param)
	if result.Error == nil {
		result.Value, result.Error = s.bizPtr.QueryWorkspace(param.Name)
	}
}

func (s *Workspace) refreshWorkspace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceListResult{}
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
	result := &imkwc.FileListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	filter := &imkwc.FileFilter{}
	filter.Parse(req)
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileList(filter.Workspace, filter)
}

func (s *Workspace) filterFileTree(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.FileTreeListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	workspace := req.URL.Query().Get("workspace")
	basePath := req.URL.Query().Get("basePath")
	stage := req.URL.Query().Get("stage")
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileTree(workspace, stage, basePath)
}

func (s *Workspace) queryFileContent(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.FileContentResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspace := req.URL.Query().Get("workspace")
	filePath := req.URL.Query().Get("filePath")
	result.Value, result.Error = s.bizPtr.QueryFileContent(workspace, filePath)
}

func (s *Workspace) chatWorkspaceContentSSE(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	if !sse.IsSSE(req) {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	workspaceName := req.URL.Query().Get("workspace")
	holder := s.sseHolderRegistry.NewHolder(res, req)
	taskFunc := func() error {
		return s.bizPtr.ChatWorkspaceContentSSE(ctx, workspaceName, holder)
	}

	err := holder.EchoSSEID()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = holder.Run(taskFunc)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Workspace) chatWorkspaceContentMessage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &imkwc.ChatMessage{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &imkwc.ChatMessageResult{}
	defer fn.PackageHTTPResponse(res, result)
	holder := s.sseHolderRegistry.GetHolder(res, req)
	if holder == nil {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	taskFunc := func() error {
		errVal := s.bizPtr.ChatWorkspaceContentMessage(ctx, param)
		if errVal != nil {
			log.Errorf("chat workspace content message error: %s", errVal.Error())
			return errors.New("chat workspace content message error")
		}

		return nil
	}
	err = holder.Run(taskFunc)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	result.Error = nil
}

func (s *Workspace) chatWorkspaceContentTerminate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &imkwc.ChatMessage{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &imkwc.ChatMessageResult{}
	defer fn.PackageHTTPResponse(res, result)
	holder := s.sseHolderRegistry.GetHolder(res, req)
	if holder == nil {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	taskFunc := func() error {
		errVal := s.bizPtr.ChatWorkspaceContentTerminate(ctx, param)
		if errVal != nil {
			log.Errorf("chat workspace content terminate error: %s", errVal.Error())
			return errors.New("chat workspace content terminate error")
		}

		return nil
	}
	err = holder.Run(taskFunc)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	result.Error = nil
}

func (s *Workspace) chatWorkspaceContentHistory(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.ChatHistoryResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceName := req.URL.Query().Get("workspace")
	result.Values, result.Error = s.bizPtr.ChatWorkspaceContentHistory(workspaceName)
}
