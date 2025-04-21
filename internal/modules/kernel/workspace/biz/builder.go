package biz

import (
	"context"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/sashabaranov/go-openai"

	"github.com/muidea/magicEngine/sse"

	imkac "github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

func (s *Workspace) constructRequest(_ *common.ChatMessageParam) *openai.ChatCompletionRequest {
	request := openai.ChatCompletionRequest{
		Model:           imkac.DeepSeekChat,
		Messages:        []openai.ChatCompletionMessage{},
		Temperature:     0.7,
		TopP:            1,
		N:               1,
		Stop:            []string{"\n\n"},
		MaxTokens:       1024,
		PresencePenalty: 0,
	}
	return &request
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context,
	param *common.ChatMessageParam,
	sink sse.Sinker) (err *cd.Error) {

	ev := event.NewEvent(imkac.ChatCompletion, s.ID(), imkac.AgentModule, nil, s.constructRequest(param))
	ev.SetData("context", ctx)
	ev.SetData("sink", sink)

	result := s.SendEvent(ev)
	err = result.Error()
	return
}

func (s *Workspace) ChatWorkspaceContentHistory(param *common.ChatHistoryParam) (ret []*common.ChatMessage, err *cd.Error) {
	return
}
