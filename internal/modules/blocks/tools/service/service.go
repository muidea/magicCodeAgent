package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Tools struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Tools
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Tools) *Tools {
	ptr := &Tools{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Tools) RegisterRoute() {
}
