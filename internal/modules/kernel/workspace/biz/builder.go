package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicEngine/sse"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type sseStream struct {
	sseID      string
	sseContext context.Context
	sseSink    sse.StreamSinker
	streamChan chan []byte
}

func newSSEStream(id string, ctx context.Context, sseSink sse.StreamSinker) *sseStream {
	return &sseStream{
		sseID:      id,
		sseContext: ctx,
		sseSink:    sseSink,
		streamChan: make(chan []byte),
	}
}

func (s *sseStream) ID() string {
	return s.sseID
}

func (s *sseStream) OnMessage(msgPtr *models.MessageInfo) {
	byteVal, _ := json.Marshal(msgPtr)
	s.streamChan <- byteVal
}

func (s *sseStream) Wait() error {
	log.Infof("workspace:%s sse stream wait", s.sseID)

	for {
		select {
		case bytes := <-s.streamChan:
			//log.Infof("workspace:%s sse stream recv：%s", s.sseID, string(bytes))
			s.sseSink.OnRecv("message", fmt.Appendf(bytes, "\n\n"))
		case <-s.sseContext.Done():
			s.sseSink.OnClose()
			log.Warnf("workspace:%s sse stream wait exit", s.sseID)
			return nil
		}
	}
}

func (s *Workspace) ChatWorkspaceContentSSE(ctx context.Context, workspaceName string, sseSinker sse.StreamSinker) (err error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		return
	}

	sseStream := newSSEStream(s.ID(), ctx, sseSinker)

	stationPtr := curWorkspace.(*Station)
	stationPtr.BindObserver(sseStream)
	defer stationPtr.UnbindObserver(sseStream)
	err = sseStream.Wait()
	return
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context, param *imkwc.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(param.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", param.Workspace))
		log.Errorf("workspace:%s not found", param.Workspace)
		return
	}
	stationPtr := curWorkspace.(*Station)

	var flowPath []*models.Flow
	flowIDVal := ctx.Value(common.FlowID{})
	if flowIDVal != nil {
		flowPath = stationPtr.GetFlowPath(flowIDVal.(string))
	} else {
		flowPath = stationPtr.GetFlowPath("")
	}

	messageList := []*models.MessageItem{}
	for _, flow := range flowPath {
		for _, message := range flow.Message {
			messageList = append(messageList, message.Send.Item)
			messageList = append(messageList, message.Recv.Item)
		}
	}

	messageItem := &models.MessageItem{
		Role:    common.User,
		Content: param.Content,
	}
	stationPtr.DisposeSendMessage(&models.MessageInfo{
		Item:       messageItem,
		CreateTime: time.Now().Local().Format(fu.CSTLayoutWithMillisecond),
		Workspace:  param.Workspace,
	})

	contentContext := context.WithValue(ctx, common.ContentSink{}, stationPtr)
	messageList = append(messageList, messageItem)
	ev := event.NewEvent(imklc.ChatCompletion, s.ID(), imklc.LLMModule, nil, messageList)
	ev.SetData("context", contentContext)
	ev.SetData("workspace", param.Workspace)
	result := s.SendEvent(ev)
	err = result.Error()
	return
}

func (s *Workspace) ChatWorkspaceContentHistory(workspaceName string) (ret []*models.MessageInfo, err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		log.Errorf("workspace:%s not found", workspaceName)
		return
	}

	stationPtr := curWorkspace.(*Station)

	messageFilter := &common.MessageFilter{
		PageIndex: 1,
		PageSize:  100,
		Workspace: workspaceName,
	}

	ret, err = stationPtr.FilterMessage(messageFilter)
	return
}
