package common

import (
	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
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
	Final  bool               `json:"final"`
	Result models.MessageItem `json:"result"`
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
