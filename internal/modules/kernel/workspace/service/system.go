package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

func (s *Workspace) querySystemConfig(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.SystemConfigResult{
		Value: &imkwc.SystemConfig{},
	}
	defer fn.PackageHTTPResponse(res, result)

	authorConfig, remoteModel, workspaceSetting := s.bizPtr.GetSystemConfig()
	result.Value.AuthorInfo = authorConfig
	result.Value.RemoteModel = remoteModel
	result.Value.WorkspaceSetting = workspaceSetting
	result.Value.VersionInfo = s.bizPtr.GetVersionInfo()
}

func (s Workspace) updateSystemConfig(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &cd.Result{}
	defer fn.PackageHTTPResponse(res, result)

	param := &imkwc.SystemConfigParam{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		result.Error = cd.NewError(cd.IllegalParam, err.Error())
		return
	}
	result.Error = s.bizPtr.SetSystemConfig(param)
}
