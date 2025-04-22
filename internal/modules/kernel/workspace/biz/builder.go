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

func (s *Workspace) constructRequest(paramPtr *common.ChatMessage) *openai.ChatCompletionRequest {
	request := openai.ChatCompletionRequest{
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
	return &request
}

func (s *Workspace) ChatWorkspaceContentMessage(ctx context.Context,
	param *common.ChatMessage,
	sink sse.Sinker) (err *cd.Error) {

	ev := event.NewEvent(imkac.ChatCompletion, s.ID(), imkac.AgentModule, nil, s.constructRequest(param))
	ev.SetData("context", ctx)
	ev.SetData("sink", sink)

	result := s.SendEvent(ev)
	err = result.Error()
	return
}

func (s *Workspace) ChatWorkspaceContentHistory(workspaceName string) (ret []*common.ChatMessage, err *cd.Error) {
	return
}
