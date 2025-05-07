package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type SwitchModeTool struct{}

func (s *SwitchModeTool) Name() string { return "switch_mode" }
func (s *SwitchModeTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *SwitchModeTool) Execute(workspaceDir string, param imbtc.ToolParam) *imbtc.ToolResult {
	modeSlug, ok := param["mode_slug"].(string)
	if !ok || modeSlug == "" {
		return &imbtc.ToolResult{Error: fmt.Errorf("missing or invalid mode identifier parameter")}
	}

	reason, _ := param["reason"].(string)

	resVal, err := SwitchMode(modeSlug, reason)
	if err != nil {
		return &imbtc.ToolResult{Error: fmt.Errorf("failed to switch mode: %v", err)}
	}

	result, err := json.MarshalIndent(resVal, "", "  ")
	if err != nil {
		return &imbtc.ToolResult{Error: fmt.Errorf("failed to serialize result: %v", err)}
	}

	return &imbtc.ToolResult{Content: result}
}

/*
## switch_mode
Description: Request to switch to a different mode. This tool allows modes to request switching to another mode when needed, such as switching to Code mode to make code changes. The user must approve the mode switch.
Parameters:
- mode_slug: (required) The slug of the mode to switch to (e.g., \"code\", \"ask\", \"architect\")
- reason: (optional) The reason for switching modes
Usage:
<switch_mode>
<mode_slug>Mode slug here</mode_slug>
<reason>Reason for switching here</reason>
</switch_mode>

Example: Requesting to switch to code mode
<switch_mode>
<mode_slug>code</mode_slug>
<reason>Need to make code changes</reason>
</switch_mode>
*/

// SwitchMode 实现switch_mode工具功能
// SwitchMode function implements the switch_mode tool functionality
func SwitchMode(modeSlug, reason string) (string, error) {
	// Return mode switching information
	if reason != "" {
		return fmt.Sprintf("Switching to mode '%s' because: %s", modeSlug, reason), nil
	}
	return fmt.Sprintf("Switching to mode '%s'", modeSlug), nil
}
