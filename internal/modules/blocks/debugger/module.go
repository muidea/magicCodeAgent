package debugger

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iirc "github.com/muidea/magicCodeAgent/internal/initators/routeregistry/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/debugger/service"
)

func init() {
	module.Register(New())
}

type Debugger struct {
	bizPtr     *biz.Debugger
	servicePtr *service.Debugger
}

func New() *Debugger {
	return &Debugger{}
}

func (s *Debugger) ID() string {
	return common.DebuggerModule
}

func (s *Debugger) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal iirc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(iirc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Debugger setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Debugger) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Debugger) Teardown() {
}
