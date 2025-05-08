package biz

import (
	"context"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
)

type Assistant interface {
	Chat(ctx context.Context, workspace string, messageItem *models.MessageItem, flowPtr *models.Flow) *cd.Error
}

func (s *Workspace) Chat(chatContext context.Context, workspace string, messageItem *models.MessageItem, flowPtr *models.Flow) (err *cd.Error) {
	flowPathList := flowPtr.GetFlowPath()
	systemMessage := &models.MessageItem{
		Content: []string{s.prompsRegistry.GetPromptsContent()},
	}

	userMessageList := models.MessageItemList{}
	for _, curPtr := range flowPathList {
		for _, message := range curPtr.Message {
			userMessageList = append(userMessageList, message.Item)
		}
	}
	userMessageList = append(userMessageList, *messageItem)
	ev := event.NewEvent(imklc.ChatCompletion, s.ID(), imklc.LLMModule, nil, userMessageList)
	ev.SetData("context", chatContext)
	ev.SetData("workspace", workspace)
	ev.SetData("systemMessage", systemMessage)
	result := s.SendEvent(ev)
	err = result.Error()
	return
}
