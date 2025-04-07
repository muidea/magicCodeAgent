package service

import (
	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicToolkit/internal/modules/kernel/monitor/biz"
)

type Monitor struct {
	roleRouteRegistry engine.RouteRegistry
	bizPtr            *biz.Monitor
}

func New(roleRouteRegistry engine.RouteRegistry, bizPtr *biz.Monitor) *Monitor {
	ptr := &Monitor{
		roleRouteRegistry: roleRouteRegistry,
		bizPtr:            bizPtr,
	}
	return ptr
}

func (s *Monitor) RegisterRoute() {
}
