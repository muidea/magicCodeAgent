package biz

import (
	"context"
	"sync"

	"github.com/sashabaranov/go-openai"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbrc "github.com/muidea/magicCodeAgent/internal/modules/blocks/recorder/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/biz/deepseek"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/biz/provider"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type LLM struct {
	biz.Base

	workspaceProvider sync.Map
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *LLM {
	ptr := &LLM{
		Base: biz.New(imklc.LLMModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(imklc.ChatCompletion, ptr.chatCompletion)
	return ptr
}

func (s *LLM) Run() (err *cd.Error) {
	return
}

func (s *LLM) constructRequest(
	systemMessage *models.MessageItem,
	userMessageList models.MessageItemList) *imklc.ChatCompletionStreamRequest {
	request := &imklc.ChatCompletionStreamRequest{
		ChatCompletionRequest: openai.ChatCompletionRequest{
			Messages:    []openai.ChatCompletionMessage{},
			Temperature: 0,   //采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
			TopP:        0.1, // 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
			N:           1,
			//Stop:            []string{"\n\n"}, //一个 string 或最多包含 16 个 string 的 list，在遇到这些词时，API 将停止生成更多的 token。
			MaxTokens:       8192, // 介于 1 到 8192 间的整数，限制一次请求中模型生成 completion 的最大 token 数。输入 token 和输出 token 的总长度受模型的上下文长度的限制。
			Stream:          true, // 如果设置为 True，将会以 SSE（server-sent events）的形式以流式发送消息增量。消息流以 data: [DONE] 结尾。
			PresencePenalty: 0,    //Default:0, 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。
		},
	}

	if systemMessage != nil {
		request.Messages = append(request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage.Content[0],
		})
	} else {
		request.Messages = append(request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		})
	}

	for _, item := range userMessageList {
		if len(item.Content) > 1 {
			chatMessage := openai.ChatCompletionMessage{
				Role: item.Role,
			}
			for _, sv := range item.Content {
				chatMessage.MultiContent = append(chatMessage.MultiContent, openai.ChatMessagePart{
					Type: openai.ChatMessagePartTypeText,
					Text: sv,
				})
			}
			request.Messages = append(request.Messages, chatMessage)
		} else {
			request.Messages = append(request.Messages, openai.ChatCompletionMessage{
				Role:    item.Role,
				Content: item.Content[0],
			})
		}
	}

	return request
}

func (s *LLM) chatCompletion(ev event.Event, re event.Result) {
	messageListVal, messageListOK := ev.Data().(models.MessageItemList)
	if !messageListOK {
		re.Set(nil, cd.NewError(cd.IllegalParam, "invalid message list"))
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

	var systemMessage *models.MessageItem
	systemMessageVal := ev.GetData("systemMessage")
	if systemMessageVal != nil {
		systemMessage = systemMessageVal.(*models.MessageItem)
	}

	requestVal := s.constructRequest(systemMessage, messageListVal)
	err := s.ChatCompletionStream(workspaceVal, contextVal, requestVal)
	if re != nil {
		if err != nil {
			log.Errorf("chat completion stream failed, err:%v", err)
			re.Set(nil, cd.NewError(cd.Unexpected, err.Error()))
			return
		}

		re.Set(nil, nil)
	}
}

func (s *LLM) ChatCompletion(
	workspaceName string,
	ctx context.Context,
	param *imklc.ChatCompletionStreamRequest) (ret *imklc.ChatCompletionStreamResponse, err error) {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.ChatCompletion(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.ChatCompletion(workspaceName, ctx, param)
}

func (s *LLM) ChatCompletionStream(
	workspaceName string,
	chatContext context.Context,
	param *imklc.ChatCompletionStreamRequest) error {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.ChatCompletionStream(workspaceName, chatContext, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.ChatCompletionStream(workspaceName, chatContext, param)
}

func (s *LLM) BetaCompletion(
	workspaceName string,
	ctx context.Context,
	param *imklc.BetaCompletionRequest) (ret *imklc.BetaCompletionStreamResponse, err error) {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.BetaCompletion(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.BetaCompletion(workspaceName, ctx, param)
}

func (s *LLM) BetaCompletionStream(
	workspaceName string,
	ctx context.Context,
	param *imklc.BetaCompletionRequest) error {
	providerVal, providerOK := s.workspaceProvider.Load(workspaceName)
	if !providerOK {
		providerPtr := provider.NewProvider(workspaceName, s)
		s.workspaceProvider.Store(workspaceName, providerPtr)
		return providerPtr.BetaCompletionStream(workspaceName, ctx, param)
	}
	providerPtr := providerVal.(*provider.Provider)
	return providerPtr.BetaCompletionStream(workspaceName, ctx, param)
}

func (s *LLM) ModelsList() (ret *imklc.ModelsListResponse, err error) {
	ret, err = deepseek.ModelsList()
	return
}

func (s *LLM) UserBalance() (ret *imklc.UserBalanceResponse, err error) {
	ret, err = deepseek.UserBalance()
	return
}

func (s *LLM) OnTerminate(msgPtr *models.MessageInfo) {
	appendEvent := event.NewEvent(imbtc.AppendMessage, s.ID(), imbtc.TracerModule, nil, msgPtr)
	appendResult := s.SendEvent(appendEvent)
	appendVal, appendErr := appendResult.Get()
	if appendErr != nil {
		log.Errorf("append message failed: %v", appendVal)
		return
	}

	debugEvent := event.NewEvent(imbrc.RecordRealtimeMessage, s.ID(), imbrc.RecorderModule, nil, appendVal)
	_ = s.SendEvent(debugEvent)
}

func (s *LLM) OnSendMessage(msgPtr *models.MessageInfo) {
	appendEvent := event.NewEvent(imbtc.AppendMessage, s.ID(), imbtc.TracerModule, nil, msgPtr)
	appendResult := s.SendEvent(appendEvent)
	appendVal, appendErr := appendResult.Get()
	if appendErr != nil {
		log.Errorf("append message failed: %v", appendVal)
		return
	}

	debugEvent := event.NewEvent(imbrc.RecordRealtimeMessage, s.ID(), imbrc.RecorderModule, nil, appendVal)
	_ = s.SendEvent(debugEvent)
}

func (s *LLM) OnRecvMessage(msgPtr *models.MessageInfo, usePtr *imklc.Usage) {
	appendEvent := event.NewEvent(imbtc.AppendMessage, s.ID(), imbtc.TracerModule, nil, msgPtr)
	appendResult := s.SendEvent(appendEvent)
	appendVal, appendErr := appendResult.Get()
	if appendErr != nil {
		log.Errorf("append message failed: %v", appendVal)
		return
	}

	debugEvent := event.NewEvent(imbrc.RecordRealtimeMessage, s.ID(), imbrc.RecorderModule, nil, appendVal)
	_ = s.SendEvent(debugEvent)
}
