package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fp "github.com/muidea/magicCommon/foundation/path"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/pkg/common"
)

type Monitor struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Monitor
}

func New(roleRouteRegistry engine.RouteRegistry, bizPtr *biz.Monitor) *Monitor {
	ptr := &Monitor{
		routeRegistry: roleRouteRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Monitor) RegisterRoute() {
	s.routeRegistry.AddHandler(common.FilterWorkSpaceFileList, engine.GET, s.filterFiles)
}

/*
| 请求参数**（Query Parameters）
| `name`       | string | 否   | 文件名模糊匹配（支持通配符 `*`，如 `*report*` 匹配包含 "report" 的文件） |
| `path`       | string | 否   | 文件路径前缀匹配（如 `/docs/2025` 匹配该路径下的所有文件）              |
| `extension`  | string | 否   | 文件扩展名精确匹配（如 `pdf`，不区分大小写）                          |
*/
func (s *Monitor) filterFiles(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.FileListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	filter := &fileFilter{}
	filter.Parse(req)
	result.Values, result.Error = s.bizPtr.FilterWorkSpaceFileList(filter)
}

type fileFilter struct {
	name      string
	path      string
	extension string
}

func (s *fileFilter) Parse(req *http.Request) {
	s.name = req.URL.Query().Get("name")
	s.path = req.URL.Query().Get("path")
	s.extension = req.URL.Query().Get("extension")
}

func (s *fileFilter) Filter(val interface{}) bool {
	filePath, fileOK := val.(string)
	if !fileOK {
		return false
	}
	if !fp.Exist(filePath) {
		return false
	}
	// TODO 其他校验条件
	return true
}
