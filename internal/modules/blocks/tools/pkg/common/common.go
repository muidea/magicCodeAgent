package common

import (
	"encoding/json"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

const ToolsModule = "/internal/modules/blocks/tools"

const (
	ExecuteTool = "/tools/execute"
)

// const ToolSectionPromptsBasePath = "assets/template/prompts/section"
const ToolsPromptsBasePath = "assets/template/prompts/tools"

type ToolParam map[string]any
type ToolResult struct {
	Content json.RawMessage `json:"content omitempty"`
	Error   error           `json:"error omitempty"`
}

type Tool interface {
	iipc.Prompt
	Execute(workspaceDir string, param ToolParam) *ToolResult
}

func ReadToolsPrompts(toolName string) (ret string, err error) {
	fileName := toolName + ".md"
	ret, err = util.GetTemplateContent(ToolsPromptsBasePath, "", fileName)
	return
}
