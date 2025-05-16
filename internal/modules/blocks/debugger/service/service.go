package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/biz"
	engine "github.com/muidea/magicEngine/http"
)

type Debugger struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Debugger
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Debugger) *Debugger {
	ptr := &Debugger{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Debugger) RegisterRoute() {
}
