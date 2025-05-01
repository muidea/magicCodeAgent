package common

import (
	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

const ToolsModule = "/internal/modules/blocks/tools"

const ToolsPromptsBasePath = "assets/template/prompts/tools"

type ToolParam map[string]any
type ToolResult []byte

type Tool interface {
	iipc.Prompts
	Execute(param ToolParam) (ToolResult, error)
}

func ReadToolsPrompts(toolName string) (ret string, err error) {
	fileName := toolName + ".md"
	ret, err = util.GetTemplateContent(ToolsPromptsBasePath, "", fileName)
	return
}
