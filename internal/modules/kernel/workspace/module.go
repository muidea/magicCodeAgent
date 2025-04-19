package workspace

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	rc "github.com/muidea/magicCodeAgent/internal/initators/routeregistry/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/service"
)

func init() {
	module.Register(New())
}

type Workspace struct {
	bizPtr     *biz.Workspace
	servicePtr *service.Workspace
}

func New() *Workspace {
	return &Workspace{}
}

func (s *Workspace) ID() string {
	return common.WorkspaceModule
}

func (s *Workspace) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal rc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(rc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Workspace setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Workspace) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Workspace) Teardown() {
}
