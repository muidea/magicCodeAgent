package knowledge

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iirc "github.com/muidea/magicCodeAgent/internal/initiators/routeregistry/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/biz"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/service"
)

func init() {
	module.Register(New())
}

type KnowledgeBase struct {
	bizPtr     *biz.KnowledgeBase
	servicePtr *service.KnowledgeBase
}

func New() *KnowledgeBase {
	return &KnowledgeBase{}
}

func (s *KnowledgeBase) ID() string {
	return imkkc.KnowledgeBaseModule
}

func (s *KnowledgeBase) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal iirc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initiator.GetEntity(iirc.RouteRegistryInitiator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Knowledge setup failed, initiator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *KnowledgeBase) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *KnowledgeBase) Teardown() {
	s.bizPtr.Stop()
}
