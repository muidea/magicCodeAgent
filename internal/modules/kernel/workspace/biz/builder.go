package biz

import (
	"context"
	"fmt"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/sashabaranov/go-openai"

	"github.com/muidea/magicEngine/sse"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	imkac "github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

type sseStream struct {
	sseID      string
	sseContext context.Context
	sseSink    sse.Sinker
	streamChan chan []byte
}

func newSSEStream(id string, ctx context.Context, sink sse.Sinker) *sseStream {
	return &sseStream{
		sseID:      id,
		sseContext: ctx,
		sseSink:    sink,
		streamChan: make(chan []byte),
	}
}

func (s *sseStream) ID() string {
	return s.sseID
}

func (s *sseStream) OnConnect(msgPtr *imkac.MessageInfo) {
	//sse.Send(msgPtr.Workspace, imbtc.OnConnect, msgPtr)
}

func (s *sseStream) OnClose(msgPtr *imkac.MessageInfo) {
	//sse.Send(msgPtr.Workspace, imbtc.OnClose, msgPtr)
}

func (s *sseStream) OnSend(msgPtr *imkac.MessageInfo) {
	//sse.Send(msgPtr.Workspace, imbtc.OnSend, msgPtr)
}

func (s *sseStream) OnRecv(msgPtr *imkac.MessageInfo) {
	s.streamChan <- []byte(msgPtr.Content)
}

func (s *sseStream) Wait() error {
	log.Infof("workspace:%s sse stream wait", s.sseID)

	for {
		select {
		case bytes := <-s.streamChan:
			log.Infof("workspace:%s sse stream recv：%s", s.sseID, string(bytes))
			s.sseSink.OnRecv(s.sseID, "data", bytes)
		case <-s.sseContext.Done():
			log.Warnf("workspace:%s sse stream wait exit", s.sseID)
			return nil
		}
	}
}

func (s *Workspace) ChatWorkspaceContentSSE(ctx context.Context, workspaceName string, sink sse.Sinker) (err error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}

	sseObserver := newSSEStream(s.ID(), ctx, sink)

	stationPtr := curWorkspace.(*Station)
	stationPtr.BindObserver(sseObserver)
	defer stationPtr.UnbindObserver(sseObserver)
	err = sseObserver.Wait()
	return
}

/*
	openai.ChatCompletionRequest{
			Model: imkac.DeepSeekChat,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: paramPtr.Content,
				},
			},
			Temperature:     0.7,
			TopP:            1,
			N:               1,
			Stop:            []string{"\n\n"},
			MaxTokens:       1024,
			PresencePenalty: 0,
		}
*/
func (s *Workspace) constructRequest(paramPtr *common.ChatMessage) *imkac.ChatCompletionRequest {
	request := &imkac.ChatCompletionRequest{
		ChatCompletionRequest: openai.ChatCompletionRequest{
			Model: imkac.DeepSeekChat,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: paramPtr.Content,
				},
			},
			Temperature:     0.7,
			TopP:            1,
			N:               1,
			Stop:            []string{"\n\n"},
			MaxTokens:       1024,
			Stream:          true,
			PresencePenalty: 0,
		},
	}

	return request
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context, param *common.ChatMessage, sink sse.Sinker) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(param.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", param.Workspace))
		log.Errorf("workspace:%s not found", param.Workspace)
		return
	}

	stationPtr := curWorkspace.(*Station)
	ev := event.NewEvent(imkac.ChatCompletion, s.ID(), imkac.AgentModule, nil, s.constructRequest(param))

	ev.SetData("workspace", param.Workspace)
	ev.SetData("context", context.WithValue(ctx, imkac.MessageSink{}, stationPtr))
	result := s.SendEvent(ev)
	err = result.Error()
	return
}

func (s *Workspace) ChatWorkspaceContentHistory(workspaceName string) (ret []*imkac.MessageInfo, err *cd.Error) {
	messageFilter := &imbtc.MessageFilter{
		PageIndex: 1,
		PageSize:  100,
		Workspace: workspaceName,
	}

	ev := event.NewEvent(imbtc.FilterMessage, s.ID(), imbtc.TracerModule, nil, messageFilter)
	resultVal, resultErr := s.SendEvent(ev).Get()
	if resultErr != nil {
		err = resultErr
		return
	}
	ret = resultVal.([]*imkac.MessageInfo)
	return
}
