package agent

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iirc "github.com/muidea/magicToolkit/internal/initators/routeregistry/pkg/common"
	"github.com/muidea/magicToolkit/internal/modules/kernel/agent/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/agent/pkg/common"
	"github.com/muidea/magicToolkit/internal/modules/kernel/agent/service"
)

func init() {
	module.Register(New())
}

type Agent struct {
	bizPtr     *biz.Agent
	servicePtr *service.Agent
}

func New() *Agent {
	return &Agent{}
}

func (s *Agent) ID() string {
	return common.AgentModule
}

func (s *Agent) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal iirc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(iirc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Agent setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Agent) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Agent) Teardown() {
}
