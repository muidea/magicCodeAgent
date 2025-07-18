package common

import (
	"fmt"
	"strings"

	"github.com/muidea/magicCodeAgent/internal/config"
	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

const ToolsModule = "/internal/modules/blocks/tools"

const (
	ReloadTools = "/tools/reload"
	ExecuteTool = "/tools/execute"
)

const innerSystemToolsPromptsBasePath = "assets/system/tools"
const innerSystemToolsSectionPromptsPath = "assets/system/prompts/tools.md"

type ToolParam map[string]any

func (s ToolParam) String() string {
	if len(s) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for k, v := range s {
		if !first {
			builder.WriteString(", ")
		}
		first = false
		builder.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	builder.WriteString("}")
	return builder.String()
}

type ToolResult struct {
	Name   string             `json:"name"`
	Result models.MessageItem `json:"result"`
}

type Tool interface {
	iipc.Prompt
	Execute(param ToolParam) (*ToolResult, error)
}

func ReadInnerToolsPrompts(toolName string) (ret string, err error) {
	fileName := toolName + ".md"
	ret, err = util.GetTemplateContent(config.GetTemplateFS(), innerSystemToolsPromptsBasePath, "", fileName)
	return
}

func ReadInnerToolsSectionPrompts() (ret string, err error) {
	ret, err = util.GetTemplateContent(config.GetTemplateFS(), innerSystemToolsSectionPromptsPath, "", "")
	return
}
