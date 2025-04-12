package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz/deepseek"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

type Agent struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Agent {
	ptr := &Agent{
		Base: biz.New(common.AgentModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Agent) Run() (err *cd.Error) {
	return
}

func (s *Agent) ChatCompletion(param *common.ChatCompletionRequest) (ret *common.ChatCompletionResponse, err error) {
	ret, err = deepseek.ChatCompletion(param)
	return
}

func (s *Agent) BetaCompletion(param *common.BetaCompletionRequest) (ret *common.BetaCompletionResponse, err error) {
	ret, err = deepseek.BetaCompletion(param)
	return
}

func (s *Agent) ModelsList() (ret *common.ModelsListResponse, err error) {
	ret, err = deepseek.ModelsList()
	return
}

func (s *Agent) UserBalance() (ret *common.UserBalanceResponse, err error) {
	ret, err = deepseek.UserBalance()
	return
}
