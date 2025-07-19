package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

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

	param := &models.Workspace{}
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

	param := &models.Workspace{}
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

	param := &models.Workspace{}
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

	param := &models.Workspace{}
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

	param := &models.Workspace{}
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
	param := &models.Workspace{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	if !models.ValidStage(param.Stage) {
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
