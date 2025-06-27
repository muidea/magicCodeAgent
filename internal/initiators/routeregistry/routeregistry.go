package routeregistry

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/initiators/routeregistry/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

func init() {
	initiator.Register(New())
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
	return common.RouteRegistryInitiator
}

func (s *routeRegistry) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	sessionRegistry := session.CreateRegistry()
	routeRegistry := engine.NewRouteRegistry()
	routeRegistry.SetApiVersion(tc.ApiVersion)
	s.eventHub = eventHub

	s.httpServer = engine.NewHTTPServer(config.BindPort(), false)
	s.httpServer.Use(engine.NewEmbedStatic(config.GetTemplateFS(), "assets", "/embed/"))
	s.httpServer.Bind(routeRegistry)

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
