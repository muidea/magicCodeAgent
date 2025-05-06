package biz

import (
	"github.com/muidea/magicCommon/event"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type Worker interface {
	DoWork(workspaceDir, content string) error
}

func (s *Workspace) DoWork(workspaceDir, content string) error {
	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	ev.SetData("workspaceDir", workspaceDir)
	s.SendEvent(ev)
	return nil
}
