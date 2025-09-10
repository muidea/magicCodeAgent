package biz

import (
	"context"
	"encoding/json"
	"fmt"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicEngine/sse"

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

	s.Base.AsyncTask(func() {
		stationPtr.BroadcastReady()
	})

	err = sseStream.Wait()
	return
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context, paramPtr *imkwc.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(paramPtr.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", paramPtr.Workspace))
		log.Errorf("workspace:%s not found", paramPtr.Workspace)
		return
	}
	stationPtr := curWorkspace.(*Station)
	messageItemPtr := &models.MessageItem{
		Role:    common.User,
		Content: []string{paramPtr.Content},
	}
	err = stationPtr.ChatMessage(ctx, messageItemPtr)
	return
}

func (s *Workspace) ChatWorkspaceContentTerminate(ctx context.Context, paramPtr *imkwc.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(paramPtr.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", paramPtr.Workspace))
		log.Errorf("workspace:%s not found", paramPtr.Workspace)
		return
	}
	stationPtr := curWorkspace.(*Station)
	err = stationPtr.ChatTerminate(ctx)
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

func (s *Workspace) FilterWorkspaceFlow(workspaceName string) (ret *models.Flow, err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(workspaceName)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", workspaceName))
		log.Errorf("workspace:%s not found", workspaceName)
		return
	}

	stationPtr := curWorkspace.(*Station)
	ret = stationPtr.GetFlows()
	return
}

func (s *Workspace) DebugWorkspaceContent(ctx context.Context, paramPtr *imkwc.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(paramPtr.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", paramPtr.Workspace))
		log.Errorf("workspace:%s not found", paramPtr.Workspace)
		return
	}
	stationPtr := curWorkspace.(*Station)
	err = stationPtr.DebugContent(ctx)
	return
}

func (s *Workspace) ChatWorkspaceTerminalContent(ctx context.Context, paramPtr *imkwc.ChatMessage) (ret string, err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(paramPtr.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", paramPtr.Workspace))
		log.Errorf("workspace:%s not found", paramPtr.Workspace)
		return
	}
	stationPtr := curWorkspace.(*Station)
	ret, err = stationPtr.ChatTerminal(ctx, paramPtr.Content)
	return
}
