package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"

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

func (s *Monitor) filterFiles(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.FileListResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	result.Values, result.Error = s.bizPtr.FilterWorkSpaceFileList(s)
}

func (s *Monitor) Filter(val interface{}) bool {
	return true
}
