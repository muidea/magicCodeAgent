package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Commander struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Commander
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Commander) *Commander {
	ptr := &Commander{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Commander) RegisterRoute() {
}
