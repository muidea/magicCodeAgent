package biz

import (
	"context"
	"sync/atomic"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz/deepseek"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

type Agent struct {
	biz.Base

	idVal int64
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Agent {
	ptr := &Agent{
		Base: biz.New(common.AgentModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(common.ChatCompletion, ptr.chatCompletion)
	return ptr
}

func (s *Agent) Run() (err *cd.Error) {
	return
}

func (s *Agent) chatCompletion(ev event.Event, re event.Result) {
	requestVal, requestOK := ev.Data().(*common.ChatCompletionRequest)
	if !requestOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid request"))
		return
	}
	workspaceVal, workspaceOK := ev.GetData("workspace").(string)
	if !workspaceOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid workspace"))
		return
	}

	contextVal, contextOK := ev.GetData("context").(context.Context)
	if !contextOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid context"))
		return
	}
	err := s.ChatCompletionStream(workspaceVal, contextVal, requestVal)
	if re != nil {
		if err != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, err.Error()))
			return
		}

		re.Set(nil, nil)
	}
}

func (s *Agent) ChatCompletion(param *common.ChatCompletionRequest) (ret *common.ChatCompletionStreamResponse, err error) {
	ret, err = deepseek.ChatCompletion(param)
	return
}

func (s *Agent) ChatCompletionStream(
	workspaceName string,
	ctx context.Context,
	param *common.ChatCompletionRequest) error {
	return deepseek.ChatCompletionStream(workspaceName, ctx, param)
}

func (s *Agent) BetaCompletion(param *common.BetaCompletionRequest) (ret *common.BetaCompletionStreamResponse, err error) {
	ret, err = deepseek.BetaCompletion(param)
	return
}

func (s *Agent) BetaCompletionStream(
	workspaceName string,
	ctx context.Context,
	param *common.BetaCompletionRequest) error {
	return deepseek.BetaCompletionStream(workspaceName, ctx, param)
}

func (s *Agent) ModelsList() (ret *common.ModelsListResponse, err error) {
	ret, err = deepseek.ModelsList()
	return
}

func (s *Agent) UserBalance() (ret *common.UserBalanceResponse, err error) {
	ret, err = deepseek.UserBalance()
	return
}

func (s *Agent) OnConnect(msgPtr *common.MessageInfo) {
	s.appendMessage(msgPtr)
}

func (s *Agent) OnClose(msgPtr *common.MessageInfo) {
	s.appendMessage(msgPtr)
}

func (s *Agent) OnSend(msgPtr *common.MessageInfo) {
	s.appendMessage(msgPtr)
}

func (s *Agent) OnRecv(msgPtr *common.MessageInfo) {
	s.appendMessage(msgPtr)
}

func (s *Agent) appendMessage(msgPtr *common.MessageInfo) {
	msgPtr.SN = atomic.AddInt64(&s.idVal, 1)
	ev := event.NewEvent(imbtc.AppendMessage, s.ID(), imbtc.TracerModule, nil, msgPtr)
	_ = s.SendEvent(ev)
}
