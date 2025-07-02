package biz

import (
	"context"

	"github.com/muidea/magicCodeAgent/internal/config"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
)

type Assistant interface {
	Chat(ctx context.Context, workspaceConfig *config.Workspace, systemPrompt, workspaceEnvPromp string, flowPathList []*models.Flow) *cd.Error
}

func (s *Workspace) Chat(
	chatContext context.Context,
	workspaceConfig *config.Workspace,
	systemPrompt string,
	workspaceEnvPrompt string,
	flowPathList []*models.Flow) (err *cd.Error) {
	if len(flowPathList) == 0 {
		err = cd.NewError(cd.InvalidParameter, "flowPathList is empty")
		return
	}

	systemMessage := &models.MessageItem{
		Content: []string{systemPrompt},
	}

	const maxCount = 30
	userMessageList := models.MessageItemList{}
	for _, curPtr := range flowPathList {
		msgList := curPtr.Message
		if len(msgList) >= maxCount {
			msgList = msgList[len(msgList)-maxCount:]
		}
		for _, message := range msgList {
			userMessageList = append(userMessageList, *message)
		}
	}
	userMessageList = append(userMessageList, models.MessageItem{
		Role:    common.User,
		Content: []string{workspaceEnvPrompt},
	})
	ev := event.NewEvent(imklc.ChatCompletion, s.ID(), imklc.LLMModule, nil, userMessageList)
	ev.SetData("context", chatContext)
	ev.SetData("workspace", workspaceConfig.Name)
	ev.SetData("systemMessage", systemMessage)
	result := s.SendEvent(ev)
	err = result.Error()
	return
}
