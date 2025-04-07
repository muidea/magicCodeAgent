package routeregistry

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/initators/routeregistry/pkg/common"
)

func init() {
	initator.Register(New())
}

func New() *routeRegistry {
	return &routeRegistry{}
}

type routeRegistry struct {
	eventHub        event.Hub
	httpServer      engine.HTTPServer
	sessionRegistry session.Registry
	routeRegistry   engine.RouteRegistry
}

func (s *routeRegistry) ID() string {
	return common.RouteRegistryInitator
}

func (s *routeRegistry) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	router := engine.NewRouteRegistry()
	sessionRegistry := session.CreateRegistry()
	routeRegistry := engine.NewRouteRegistry()

	s.eventHub = eventHub

	s.httpServer = engine.NewHTTPServer(config.BindPort())
	s.httpServer.Bind(router)

	s.routeRegistry = routeRegistry
	s.sessionRegistry = sessionRegistry
	return
}

func (s *routeRegistry) Run() {
	go func() {
		s.httpServer.Run()
	}()
}

func (s *routeRegistry) Teardown() {
}

func (s *routeRegistry) GetRouteRegistry() (ret engine.RouteRegistry) {
	ret = s.routeRegistry
	return
}
