package service

import (
	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/biz"
	imbcc "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/pkg/common"
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
	s.routeRegistry.AddHandler(imbcc.ConsoleLogsAppend, engine.POST, s.ConsoleLogs)
}
