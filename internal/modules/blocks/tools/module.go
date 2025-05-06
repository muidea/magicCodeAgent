package tools

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"

	_ "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/config"
)

func init() {
	module.Register(New())
}

type Tools struct {
	bizPtr *biz.Tools
}

func New() *Tools {
	return &Tools{}
}

func (s *Tools) ID() string {
	return common.ToolsModule
}

func (s *Tools) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var promptsRegistryHelperVal iipc.PromptsRegistryHeper
	var promptsRegistryHelperErr *cd.Error
	promptsRegistryHelperVal, promptsRegistryHelperErr = initator.GetEntity(iipc.PromptsInitator, promptsRegistryHelperVal)
	if promptsRegistryHelperErr != nil {
		err = promptsRegistryHelperErr
		log.Errorf("Tool setup failed, initator.GetEntity error:%s", promptsRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine, promptsRegistryHelperVal.GetPromptsRegistry())
	return nil
}

func (s *Tools) Run() (err *cd.Error) {
	err = s.bizPtr.Run()
	if err != nil {
		return
	}
	return
}

func (s *Tools) Teardown() {
}
