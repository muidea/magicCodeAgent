package biz

import (
	"context"
	"sync/atomic"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicEngine/sse"

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
	contextVal, contextOK := ev.GetData("context").(context.Context)
	if !contextOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid context"))
		return
	}
	sinkVal, sinkOK := ev.GetData("sink").(sse.Sinker)
	if !sinkOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid sink"))
		return
	}
	err := s.ChatCompletionStream(contextVal, requestVal, sinkVal)
	if re != nil {
		re.Set(nil, cd.NewError(cd.Unexpected, err.Error()))
	}
}

func (s *Agent) ChatCompletion(param *common.ChatCompletionRequest) (ret *common.ChatCompletionStreamResponse, err error) {
	ret, err = deepseek.ChatCompletion(param)
	return
}

func (s *Agent) ChatCompletionStream(
	ctx context.Context,
	param *common.ChatCompletionRequest,
	sink sse.Sinker) error {
	return deepseek.ChatCompletionStream(ctx, param, sink, s)
}

func (s *Agent) BetaCompletion(param *common.BetaCompletionRequest) (ret *common.BetaCompletionStreamResponse, err error) {
	ret, err = deepseek.BetaCompletion(param)
	return
}

func (s *Agent) BetaCompletionStream(ctx context.Context, param *common.BetaCompletionRequest, sink sse.Sinker) error {
	return deepseek.BetaCompletionStream(ctx, param, sink, s)
}

func (s *Agent) ModelsList() (ret *common.ModelsListResponse, err error) {
	ret, err = deepseek.ModelsList()
	return
}

func (s *Agent) ModelsListStream(ctx context.Context, sink sse.Sinker) error {
	return deepseek.ModelsListStream(ctx, sink, s)
}

func (s *Agent) UserBalance() (ret *common.UserBalanceResponse, err error) {
	ret, err = deepseek.UserBalance()
	return
}

func (s *Agent) UserBalanceStream(ctx context.Context, sink sse.Sinker) error {
	return deepseek.UserBalanceStream(ctx, sink, s)
}

func (s *Agent) AppendMessage(msg *common.MessageInfo) {
	msg.SN = atomic.AddInt64(&s.idVal, 1)

	ev := event.NewEvent(imbtc.AppendMessage, s.ID(), imbtc.TracerModule, nil, msg)
	_ = s.SendEvent(ev)
}
