package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledge/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Knowledge struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Knowledge
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Knowledge) *Knowledge {
	ptr := &Knowledge{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Knowledge) RegisterRoute() {
}
