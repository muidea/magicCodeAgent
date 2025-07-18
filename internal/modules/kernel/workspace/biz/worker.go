package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Worker interface {
	DoWork(workspaceConfig *models.Workspace, content string) (param imbtc.ToolParam, result *imbtc.ToolResult, err *cd.Error)
}

func (s *Workspace) DoWork(workspaceConfig *models.Workspace, content string) (param imbtc.ToolParam, result *imbtc.ToolResult, err *cd.Error) {
	executeEvent := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	executeResult := s.SendEvent(executeEvent)
	resultVal, resultErr := executeResult.Get()
	if resultErr != nil {
		log.Errorf("do work failed, %s", resultErr)
		err = resultErr
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	toolResultPtr, toolResultOK := resultVal.(*imbtc.ToolResult)
	if !toolResultOK {
		err = cd.NewError(cd.Unexpected, "do work failed, illegal result")
		return //"", fmt.Errorf("do work failed, %s", resultErr)
	}

	param = executeResult.GetVal("param").(imbtc.ToolParam)
	result = toolResultPtr
	return
}
