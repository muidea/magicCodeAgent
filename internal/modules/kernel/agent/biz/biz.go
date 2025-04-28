package biz

import (
	"context"
	"sync"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz/deepseek"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz/provider"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Agent struct {
	biz.Base

	workspaceProvider sync.Map
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

func (s *Agent) ChatCompletion(
	workspaceName string,
	ctx context.Context,
	param *common.ChatCompletionRequest) (ret *common.ChatCompletionStreamResponse, err error) {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.ChatCompletion(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.ChatCompletion(workspaceName, ctx, param)
}

func (s *Agent) ChatCompletionStream(
	workspaceName string,
	chatContext context.Context,
	param *common.ChatCompletionRequest) error {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.ChatCompletionStream(workspaceName, chatContext, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.ChatCompletionStream(workspaceName, chatContext, param)
}

func (s *Agent) BetaCompletion(
	workspaceName string,
	ctx context.Context,
	param *common.BetaCompletionRequest) (ret *common.BetaCompletionStreamResponse, err error) {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.BetaCompletion(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.BetaCompletion(workspaceName, ctx, param)
}

func (s *Agent) BetaCompletionStream(
	workspaceName string,
	ctx context.Context,
	param *common.BetaCompletionRequest) error {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.BetaCompletionStream(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.BetaCompletionStream(workspaceName, ctx, param)
}

func (s *Agent) ModelsList() (ret *common.ModelsListResponse, err error) {
	ret, err = deepseek.ModelsList()
	return
}

func (s *Agent) UserBalance() (ret *common.UserBalanceResponse, err error) {
	ret, err = deepseek.UserBalance()
	return
}

func (s *Agent) OnMessage(msgPtr *models.MessageInfo) {
	s.broadCastMessage(msgPtr)
}

func (s *Agent) broadCastMessage(msgPtr *models.MessageInfo) {
	ev := event.NewEvent(common.BroadcastMessageInfo, s.ID(), s.RootDestination(), nil, msgPtr)
	s.PostEvent(ev)
}
