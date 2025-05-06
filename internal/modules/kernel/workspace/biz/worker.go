package biz

import (
	"github.com/muidea/magicCommon/event"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type Worker interface {
	DoWork(content string) error
}

func (s *Workspace) DoWork(workspace, content string) error {
	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, content)
	ev.SetData("workspace", workspace)
	s.SendEvent(ev)
	return nil
}
