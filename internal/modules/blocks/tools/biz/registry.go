package biz

import (
	"fmt"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type ToolRegistry struct {
	promptsRegistry iipc.PromptsRegistry
	name2Tools      map[string]imbtc.Tool
}

func NewToolRegistry(promptsRegistry iipc.PromptsRegistry) *ToolRegistry {
	ptr := &ToolRegistry{
		promptsRegistry: promptsRegistry,
		name2Tools:      make(map[string]imbtc.Tool),
	}

	if promptsRegistry != nil {
		promptsRegistry.AddSection(ptr)
	}

	return ptr
}

func (s *ToolRegistry) Name() string {
	return "Tools"
}

func (s *ToolRegistry) Prompt() string {
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", "tools.md")
	if contentErr == nil {
		return contentVal
	}

	log.Warnf("failed to parse tool prompts, error:%s", contentErr.Error())
	return ""
}

func (s *ToolRegistry) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ToolRegistry) Register(tool imbtc.Tool) {
	if s.promptsRegistry != nil {
		s.promptsRegistry.AddPrompts(s.Name(), tool)
	}
	s.name2Tools[tool.Name()] = tool
}

func (s *ToolRegistry) Execute(workspaceDir, stage, toolName string, param imbtc.ToolParam) *imbtc.ToolResult {
	log.Infof("execute tool, workspaceDir:%s, stage:%s, toolName:%s, param:%v", workspaceDir, stage, toolName, param)
	tool, ok := s.name2Tools[toolName]
	if !ok {
		// 获取所有可用工具名称列表
		availableTools := make([]string, 0, len(s.name2Tools))
		for name := range s.name2Tools {
			availableTools = append(availableTools, name)
		}
		return &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("tool '%s' does not exist. Available tools: %v", toolName, availableTools),
				},
			},
		}
	}

	return tool.Execute(workspaceDir, stage, param)
}
