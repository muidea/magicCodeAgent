package common

import (
	"fmt"
	"strings"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

const ToolsModule = "/internal/modules/blocks/tools"

const (
	ExecuteTool = "/tools/execute"
)

// const ToolSectionPromptsBasePath = "assets/prompts/section"
const ToolsPromptsBasePath = "assets/prompts/tools"

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
	Final  bool               `json:"final"`
	Result models.MessageItem `json:"result"`
}

type Tool interface {
	iipc.Prompt
	Execute(workspaceName, workspaceDir, stage string, param ToolParam) *ToolResult
}

func ReadToolsPrompts(toolName string) (ret string, err error) {
	fileName := toolName + ".md"
	ret, err = util.GetTemplateContent(ToolsPromptsBasePath, "", fileName)
	return
}
