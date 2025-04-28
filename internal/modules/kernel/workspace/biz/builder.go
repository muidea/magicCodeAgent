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

	imkac "github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
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
			s.sseWriter.SendData("message", fmt.Appendf(bytes, "\n\n"))
		case <-s.sseContext.Done():
			log.Warnf("workspace:%s sse stream wait exit", s.sseID)
			return nil
		}
	}
}

func (s *Workspace) recvMessage(ev event.Event, re event.Result) {
	msgPtr, ok := ev.Data().(*models.MessageInfo)
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
func (s *Workspace) constructRequest(paramPtr *imkwc.ChatMessage) *imkac.ChatCompletionRequest {
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
			Temperature:     0,   //采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
			TopP:            0.1, // 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
			N:               1,
			Stop:            []string{"\n\n"}, //一个 string 或最多包含 16 个 string 的 list，在遇到这些词时，API 将停止生成更多的 token。
			MaxTokens:       8192,             // 介于 1 到 8192 间的整数，限制一次请求中模型生成 completion 的最大 token 数。输入 token 和输出 token 的总长度受模型的上下文长度的限制。
			Stream:          true,             // 如果设置为 True，将会以 SSE（server-sent events）的形式以流式发送消息增量。消息流以 data: [DONE] 结尾。
			PresencePenalty: 0,                //Default:0, 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。
		},
	}

	return request
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context, param *imkwc.ChatMessage) (err *cd.Error) {
	curWorkspace, ok := s.workspaceStationCache.Load(param.Workspace)
	if !ok {
		err = cd.NewError(cd.NotFound, fmt.Sprintf("workspace:%s not found", param.Workspace))
		log.Errorf("workspace:%s not found", param.Workspace)
		return
	}

	stationPtr := curWorkspace.(*Station)
	stationPtr.DisposeSendMessage(&models.MessageInfo{
		Type:       common.SendMsg,
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
