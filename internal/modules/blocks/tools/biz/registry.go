package biz

import (
	"errors"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
	"github.com/muidea/magicCommon/foundation/log"
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

	log.Warnf("parse tool prompts failed, err:%s", contentErr.Error())
	return ""
}

func (s *ToolRegistry) Register(tool imbtc.Tool) {
	if s.promptsRegistry != nil {
		s.promptsRegistry.AddPrompts(s.Name(), tool)
	}
	s.name2Tools[tool.Name()] = tool
}

func (s *ToolRegistry) Execute(workspaceDir, toolName string, param imbtc.ToolParam) (imbtc.ToolResult, error) {
	tool, ok := s.name2Tools[toolName]
	if !ok {
		return nil, errors.New("not exist tool")
	}

	return tool.Execute(workspaceDir, param)
}
