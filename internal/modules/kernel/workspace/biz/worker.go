package biz

import (
	"path"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Worker interface {
	DoWork(workspaceDir, workMode, content string) (result models.MessageItem, final bool)
}

func (s *Workspace) DoWork(workspaceDir, workMode, content string) (result models.MessageItem, final bool) {
	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	ev.SetData("workspaceDir", path.Join(workspaceDir, workMode))
	resultVal, resultErr := s.SendEvent(ev).Get()
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
