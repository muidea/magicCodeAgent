package biz

import (
	"fmt"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type ToolRegistry struct {
	currentWorkspace *models.Workspace
	commander        Commander
	templateQuerier  TemplateQuerier
	promptsRegistry  iipc.PromptsRegistry
	name2Tools       map[string]imbtc.Tool
}

func NewToolRegistry(commander Commander, templateQuerier TemplateQuerier, promptsRegistry iipc.PromptsRegistry) *ToolRegistry {
	ptr := &ToolRegistry{
		commander:       commander,
		templateQuerier: templateQuerier,
		promptsRegistry: promptsRegistry,
		name2Tools:      make(map[string]imbtc.Tool),
	}

	return ptr
}

func (s *ToolRegistry) Reload(workspacePtr *models.Workspace) (err *cd.Error) {
	if s.currentWorkspace != nil && s.currentWorkspace.IsSamge(workspacePtr) {
		return
	}

	templatePtr, templateErr := s.templateQuerier.QueryTemplate(workspacePtr.Template)
	if templateErr != nil {
		log.Errorf("loading %s template failed, error: %v", workspacePtr.Name, templateErr.Error())
		return
	}

	s.currentWorkspace = workspacePtr
	s.promptsRegistry.Reload(workspacePtr)
	s.promptsRegistry.AddSection(s)
	s.name2Tools = make(map[string]imbtc.Tool)
	s.loadInnerTools(templatePtr)
	s.loadWorkspaceTools(templatePtr)
	return
}

func (s *ToolRegistry) loadInnerTools(templatePtr *models.Template) {
	s.Register(&ReadFileTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ReadStageFileTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&FetchInstructionsTool{})
	s.Register(&SearchFilesTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ListFilesTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ListStageFilesTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ListCodeDefinitionNamesTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ApplyDiffTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&WriteToFileTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&InsertContentTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&SearchAndReplaceTool{
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&ExecuteCommandTool{
		commander:    s.commander,
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&AttemptCompletionTool{
		commander:    s.commander,
		workspacePtr: s.currentWorkspace,
	})
	s.Register(&AttemptFinishTool{})
	s.Register(&DecomposeDesignPlanTool{
		workspacePtr: s.currentWorkspace,
		templatePtr:  templatePtr,
	})
	s.Register(&DecomposeCodePlanTool{
		workspacePtr: s.currentWorkspace,
		templatePtr:  templatePtr,
	})
	s.Register(&DecomposeDebugPlanTool{
		workspacePtr: s.currentWorkspace,
		templatePtr:  templatePtr,
	})
}

func (s *ToolRegistry) loadWorkspaceTools(templatePtr *models.Template) {
	toolsList := templatePtr.LoadTools()
	for _, val := range toolsList {
		toolPtr, toolErr := NewWorkspaceTool(s.commander, s.currentWorkspace, val)
		if toolErr != nil {
			log.Errorf("loading %s tool failed, error: %v", s.currentWorkspace.Name, toolErr.Error())
			continue
		}

		s.Register(toolPtr)
	}
}

func (s *ToolRegistry) Name() string {
	return "Tools"
}

func (s *ToolRegistry) Prompt() string {
	contentVal, contentErr := imbtc.ReadInnerToolsSectionPrompts()
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

func (s *ToolRegistry) Execute(toolName string, param imbtc.ToolParam) (*imbtc.ToolResult, error) {
	//log.Infof("execute tool, workspaceName:%s, workspaceDir:%s, stage:%s, toolName:%s, param:%v", workspaceName, workspaceDir, stage, toolName, param)
	tool, ok := s.name2Tools[toolName]
	if !ok {
		// 获取所有可用工具名称列表
		availableTools := make([]string, 0, len(s.name2Tools))
		for name := range s.name2Tools {
			availableTools = append(availableTools, name)
		}
		return &imbtc.ToolResult{
			Name: toolName,
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("tool '%s' does not exist. Available tools: %v", toolName, availableTools),
				},
			},
		}, fmt.Errorf("tool '%s' does not exist", toolName)
	}

	return tool.Execute(param)
}
