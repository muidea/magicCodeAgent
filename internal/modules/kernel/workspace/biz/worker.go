package biz

import (
	"github.com/muidea/magicCommon/event"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type Worker interface {
	DoWork(workspaceDir, content string) (string, error)
}

func (s *Workspace) DoWork(workspaceDir, content string) (string, error) {
	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	ev.SetData("workspaceDir", workspaceDir)
	result := s.SendEvent(ev)
	resultVal, resultErr := result.Get()
	if resultErr != nil {
		return "", resultErr
	}

	return resultVal.(string), nil
}
