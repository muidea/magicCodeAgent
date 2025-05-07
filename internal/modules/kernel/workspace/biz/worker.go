package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type Worker interface {
	DoWork(workspaceDir, content string) (result string, final bool)
}

func (s *Workspace) DoWork(workspaceDir, content string) (result string, final bool) {
	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	ev.SetData("workspaceDir", workspaceDir)
	resultVal, resultErr := s.SendEvent(ev).Get()
	if resultErr != nil {
		result = resultErr.Error()
		final = true
		log.Errorf("do work failed, %s", resultErr)
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	toolResultPtr, toolResultOK := resultVal.(*imbtc.ToolResult)
	if !toolResultOK {
		result = "do work failed, illegal result"
		final = true
		log.Errorf(result)
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	if toolResultPtr.Error == nil && toolResultPtr.Content == nil {
		final = toolResultPtr.Final
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	result = string(toolResultPtr.Content)
	final = toolResultPtr.Final
	return
}
