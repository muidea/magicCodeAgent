package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicEngine/sse"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	imkac "github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

type sseStream struct {
	sseID      string
	sseContext context.Context
	sseWriter  sse.Writer
	streamChan chan []byte
}

func newSSEStream(id string, ctx context.Context, sseWriter sse.Writer) *sseStream {
	return &sseStream{
		sseID:      id,
		sseContext: ctx,
		sseWriter:  sseWriter,
		streamChan: make(chan []byte),
	}
}

func (s *sseStream) ID() string {
	return s.sseID
}

func (s *sseStream) OnMessage(msgPtr *imkac.MessageInfo) {
	byteVal, _ := json.Marshal(msgPtr)
	s.streamChan <- byteVal
}

func (s *sseStream) Wait() error {
	log.Infof("workspace:%s sse stream wait", s.sseID)

	for {
		select {
		case bytes := <-s.streamChan:
			log.Infof("workspace:%s sse stream recv：%s", s.sseID, string(bytes))
			s.sseWriter.SendData("message", fmt.Appendf(bytes, "\n\n"))
		case <-s.sseContext.Done():
			log.Warnf("workspace:%s sse stream wait exit", s.sseID)
			return nil
		}
	}
}

func (s *Workspace) recvMessage(ev event.Event, re event.Result) {
	msgPtr, ok := ev.Data().(*imkac.MessageInfo)
	if !ok {
		log.Errorf("workspace:%s recv message failed, invalid message:%v")
		return
	}
	stationVal, stationOK := s.workspaceStationCache.Load(msgPtr.Workspace)
	if !stationOK {
		log.Errorf("workspace:%s recv message failed, station not found", msgPtr.Workspace)
		return
	}

	stationVal.(*Station).DisposeRecvMessage(msgPtr)
	if re != nil {
		re.Set(nil, nil)
	}
}

func (s *Workspace) ChatWorkspaceContentSSE(ctx context.Context, workspaceName string, sseWriter sse.Writer) (err error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}

	sseStream := newSSEStream(s.ID(), ctx, sseWriter)

	stationPtr := curWorkspace.(*Station)
	stationPtr.BindObserver(sseStream)
	defer stationPtr.UnbindObserver(sseStream)
	err = sseStream.Wait()
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

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context, param *common.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(param.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", param.Workspace))
		log.Errorf("workspace:%s not found", param.Workspace)
		return
	}

	stationPtr := curWorkspace.(*Station)
	stationPtr.DisposeSendMessage(&imkac.MessageInfo{
		Type:       imkac.SendMsg,
		Content:    param.Content,
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  param.Workspace,
	})

	ev := event.NewEvent(imkac.ChatCompletion, s.ID(), imkac.AgentModule, nil, s.constructRequest(param))
	ev.SetData("context", ctx)
	ev.SetData("workspace", param.Workspace)
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
