package biz

import (
	"context"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
)

type Assistant interface {
	Chat(ctx context.Context, workspace string, messageItem *models.MessageItem, flowPtr *models.Flow) *cd.Error
}

func (s *Workspace) Chat(chatContext context.Context, workspace string, messageItem *models.MessageItem, flowPtr *models.Flow) (err *cd.Error) {
	flowPath := flowPtr.GetFlowPath()
	systemMessage := &models.MessageItem{
		Role:    common.System,
		Content: s.prompsRegistry.GetPromptsContent(),
	}

	userMessageList := []*models.MessageItem{}
	for _, flow := range flowPath {
		for _, message := range flow.Message {
			if message.Send != nil {
				userMessageList = append(userMessageList, message.Send.Item)
			}
			if message.Recv != nil {
				userMessageList = append(userMessageList, message.Recv.Item)
			}
		}
	}

	log.Infof("%v", messageItem.Content)

	userMessageList = append(userMessageList, messageItem)
	ev := event.NewEvent(imklc.ChatCompletion, s.ID(), imklc.LLMModule, nil, userMessageList)
	ev.SetData("context", chatContext)
	ev.SetData("workspace", workspace)
	ev.SetData("systemMessage", systemMessage)
	result := s.SendEvent(ev)
	err = result.Error()
	return
}
