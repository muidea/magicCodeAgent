package service

import (
	"context"
	"net/http"

	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	fn "github.com/muidea/magicCommon/foundation/net"
)

func (s Workspace) queryWorkspaceTemplate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.WorkspaceTemplateResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceType := req.URL.Query().Get("workspaceType")
	result.Value, result.Error = s.bizPtr.QueryWorkspaceTemplate(workspaceType)
}
