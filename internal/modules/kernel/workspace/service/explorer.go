package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

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
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFileTree(workspace, basePath)
}

func (s *Workspace) queryFileContent(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.FileContentResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspace := req.URL.Query().Get("workspace")
	filePath := req.URL.Query().Get("filePath")
	result.Value, result.Error = s.bizPtr.QueryFileContent(workspace, filePath)
}
