package biz

import (
	"context"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
)

type Assistant interface {
	Chat(ctx context.Context, workspace string, flowPathList []*models.Flow) *cd.Error
}

func (s *Workspace) Chat(chatContext context.Context, workspace string, flowPathList []*models.Flow) (err *cd.Error) {
	if len(flowPathList) == 0 {
		err = cd.NewError(cd.InvalidParameter, "flowPathList is empty")
		return
	}

	systemMessage := &models.MessageItem{
		Content: []string{s.prompsRegistry.GetPromptsContent()},
	}

	userMessageList := models.MessageItemList{}
	for _, curPtr := range flowPathList {
		for _, message := range curPtr.Message {
			userMessageList = append(userMessageList, message.Item)
		}
	}
	ev := event.NewEvent(imklc.ChatCompletion, s.ID(), imklc.LLMModule, nil, userMessageList)
	ev.SetData("context", chatContext)
	ev.SetData("workspace", workspace)
	ev.SetData("systemMessage", systemMessage)
	result := s.SendEvent(ev)
	err = result.Error()
	return
}
