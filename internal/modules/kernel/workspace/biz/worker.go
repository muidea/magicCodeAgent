package biz

import (
	"fmt"
	"os"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	"github.com/muidea/magicCodeAgent/internal/config"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Worker interface {
	DoWork(workspaceConfig *config.Workspace, content string) (result models.MessageItem, final bool)
}

func (s *Workspace) DoWork(workspaceConfig *config.Workspace, content string) (result models.MessageItem, final bool) {
	workspaceDir := workspaceConfig.GetWorkspaceDir()
	if !fp.Exist(workspaceConfig.GetWorkspaceDir()) {
		if err := os.MkdirAll(workspaceDir, 0755); err != nil {
			result = models.MessageItem{
				Content: []string{fmt.Sprintf("create work dir failed, error%s", err.Error())},
			}
			final = true
			return
		}
	}
	executeEvent := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	executeEvent.SetData("workspaceName", workspaceConfig.GetName())
	executeEvent.SetData("workspaceDir", workspaceDir)
	executeEvent.SetData("stage", workspaceConfig.GetStage())
	executeResult := s.SendEvent(executeEvent)
	resultVal, resultErr := executeResult.Get()
	if resultErr != nil {
		result = models.MessageItem{
			Content: []string{resultErr.Error()},
		}
		final = true
		log.Errorf("do work failed, %s", resultErr)
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	toolResultPtr, toolResultOK := resultVal.(*imbtc.ToolResult)
	if !toolResultOK {
		result = models.MessageItem{
			Content: []string{"do work failed, illegal result"},
		}
		final = true
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	result = toolResultPtr.Result
	final = toolResultPtr.Final
	return
}
