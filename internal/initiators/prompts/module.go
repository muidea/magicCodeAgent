package prompts

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/initiator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
)

func init() {
	initiator.Register(New())
}

type Prompts struct {
	promptsRegistry *promptsRegistry
}

func New() *Prompts {
	return &Prompts{
		promptsRegistry: NewPromptsRegistry(),
	}
}

func (s *Prompts) ID() string {
	return common.PromptsInitiator
}

func (s *Prompts) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	return nil
}

func (s *Prompts) Run() (err *cd.Error) {
	return
}

func (s *Prompts) Teardown() {
}

func (s *Prompts) GetPromptsRegistry() common.PromptsRegistry {
	return s.promptsRegistry
}
