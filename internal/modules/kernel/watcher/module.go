package watcher

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	rc "github.com/muidea/magicToolkit/internal/initators/routeregistry/pkg/common"
	"github.com/muidea/magicToolkit/internal/modules/kernel/watcher/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/watcher/pkg/common"
	"github.com/muidea/magicToolkit/internal/modules/kernel/watcher/service"
)

func init() {
	module.Register(New())
}

type Watcher struct {
	bizPtr     *biz.Watcher
	servicePtr *service.Watcher
}

func New() *Watcher {
	return &Watcher{}
}

func (s *Watcher) ID() string {
	return common.WatcherModule
}

func (s *Watcher) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal rc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(rc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Watcher setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Watcher) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Watcher) Teardown() {
}
