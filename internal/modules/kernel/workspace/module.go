package workspace

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iiec "github.com/muidea/magicCodeAgent/internal/initiators/embedtext/pkg/common"
	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	iirc "github.com/muidea/magicCodeAgent/internal/initiators/routeregistry/pkg/common"
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
	var routeRegistryHelperVal iirc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initiator.GetEntity(iirc.RouteRegistryInitiator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("Workspace setup failed, initiator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	var promptsRegistryHelperVal iipc.PromptsRegistryHelper
	var promptsRegistryHelperErr *cd.Error
	promptsRegistryHelperVal, promptsRegistryHelperErr = initiator.GetEntity(iipc.PromptsInitiator, promptsRegistryHelperVal)
	if promptsRegistryHelperErr != nil {
		err = promptsRegistryHelperErr
		log.Errorf("Workspace setup failed, initiator.GetEntity error:%s", promptsRegistryHelperErr.Error())
		return
	}

	var embedTextWorkerHelperVal iiec.EmbedTextHelper
	var embedTextWorkerHelperErr *cd.Error
	embedTextWorkerHelperVal, embedTextWorkerHelperErr = initiator.GetEntity(iiec.EmbedtextInitiator, embedTextWorkerHelperVal)
	if embedTextWorkerHelperErr != nil {
		err = embedTextWorkerHelperErr
		log.Errorf("Knowledge setup failed, initiator.GetEntity error:%s", embedTextWorkerHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine, promptsRegistryHelperVal.GetPromptsRegistry(), embedTextWorkerHelperVal.GetEmbedText())

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
