package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tool/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Tool struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Tool
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Tool) *Tool {
	ptr := &Tool{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Tool) RegisterRoute() {
}
