package service

import (
	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicToolkit/internal/modules/kernel/watcher/biz"
)

type Watcher struct {
	roleRouteRegistry engine.RouteRegistry
	bizPtr            *biz.Watcher
}

func New(roleRouteRegistry engine.RouteRegistry, bizPtr *biz.Watcher) *Watcher {
	ptr := &Watcher{
		roleRouteRegistry: roleRouteRegistry,
		bizPtr:            bizPtr,
	}
	return ptr
}

func (s *Watcher) RegisterRoute() {
}
