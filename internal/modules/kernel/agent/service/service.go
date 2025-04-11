package service

import (
	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicToolkit/internal/modules/kernel/agent/biz"
)

type Agent struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Agent
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Agent) *Agent {
	ptr := &Agent{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Agent) RegisterRoute() {
}
