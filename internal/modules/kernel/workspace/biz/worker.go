package biz

import (
	"errors"
	"strings"

	"github.com/muidea/magicCommon/event"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type Worker interface {
	DoWork(content string) error
}

func (s *Workspace) DoWork(workspace, content string) error {
	parts := strings.Split(content, "\n\n")
	if len(parts) < 2 {
		return errors.New("invalid response format")
	}

	hasThinking := false
	hasTool := false

	var toolsContent string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if strings.HasPrefix(trimmed, "<thinking>") && strings.HasSuffix(trimmed, "</thinking>") {
			hasThinking = true
		} else if strings.HasPrefix(trimmed, "<tool>") && strings.HasSuffix(trimmed, "</tool>") {
			hasTool = true
			toolsContent = trimmed
		}
		if hasThinking && hasTool {
			break // 提前退出循环
		}
	}

	if !hasThinking || !hasTool {
		return errors.New("invalid response format")
	}

	ev := event.NewEvent(imbtc.ExecuteTool, s.ID(), imbtc.ToolsModule, nil, toolsContent)
	ev.SetData("workspace", workspace)
	s.SendEvent(ev)
	return nil
}
