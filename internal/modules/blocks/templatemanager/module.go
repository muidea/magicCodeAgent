package templatemanager

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iirc "github.com/muidea/magicCodeAgent/internal/initiators/routeregistry/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/service"
)

func init() {
	module.Register(New())
}

type Templatemanager struct {
	bizPtr     *biz.Templatemanager
	servicePtr *service.Templatemanager
}

func New() *Templatemanager {
	return &Templatemanager{}
}

func (s *Templatemanager) ID() string {
	return common.TemplatemanagerModule
}

func (s *Templatemanager) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal iirc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initiator.GetEntity(iirc.RouteRegistryInitiator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Templatemanager setup failed, initiator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *Templatemanager) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *Templatemanager) Teardown() {
}
