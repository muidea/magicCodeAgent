package monitor

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	rc "github.com/muidea/magicCodeAgent/internal/initators/routeregistry/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/monitor/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/monitor/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/monitor/service"
)

func init() {
	module.Register(New())
}

type Monitor struct {
	bizPtr     *biz.Monitor
	servicePtr *service.Monitor
}

func New() *Monitor {
	return &Monitor{}
}

func (s *Monitor) ID() string {
	return common.MonitorModule
}

func (s *Monitor) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal rc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(rc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Monitor setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Monitor) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Monitor) Teardown() {
}
