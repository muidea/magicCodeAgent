package biz

import (
	"fmt"
	"path/filepath"
	"strings"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type DecomposeDesignSubtaskTask struct {
	promptsRegistry iipc.PromptsRegistry
}

func (s *DecomposeDesignSubtaskTask) Name() string {
	return "decompose_design_subtask"
}

func (s *DecomposeDesignSubtaskTask) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeDesignSubtaskTask) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeDesignSubtaskTask) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s, param:%v", s.Name(), param),
				},
			},
		}
		if err != nil {
			ret.Result.Content = append(ret.Result.Content, err.Error())
		} else {
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	name, ok := param["name"].(string)
	if !ok || name == "" {
		err = fmt.Errorf("name is empty")
		return
	}
	level, ok := param["level"].(string)
	if !ok || level == "" {
		err = fmt.Errorf("level is empty")
		return
	}
	description, ok := param["description"].(string)
	if !ok || description == "" {
		err = fmt.Errorf("description is empty")
		return
	}

	stringBuffer := strings.Builder{}
	stringBuffer.WriteString("## " + name + "\n\n")
	stringBuffer.WriteString(description + "\n\n")

	levelSectionName := fmt.Sprintf("decompose-design-%s-plan", level)
	content := s.promptsRegistry.ReadSectionPromptContent(levelSectionName)
	if content != "" {
		stringBuffer.WriteString(content)
		stringBuffer.WriteString("\n\n")
	}

	var finalPath string
	if level != "project" {
		// 查看兄弟层级分解清单
		finalPath = filepath.Join(workspaceDir, stage)
		fileList, fileErr := fp.ListPath(finalPath, "*.md", false)
		if fileErr != nil {
			err = fileErr
			return
		}
		if len(fileList) > 0 {
			stringBuffer.WriteString("## 同级任务清单\n\n")
			for _, item := range fileList {
				stringBuffer.WriteString("- " + item[len(finalPath)+1:] + "\n")
			}
			stringBuffer.WriteString("\n")
		}
	}

	// 查看子层级分解清单
	finalPath = filepath.Join(workspaceDir, stage, name)
	fileList, fileErr := fp.ListPath(finalPath, "*.md", false)
	if fileErr != nil {
		err = fileErr
		return
	}
	if len(fileList) > 0 {
		stringBuffer.WriteString("## 子集任务清单\n\n")
		for _, item := range fileList {
			stringBuffer.WriteString("- " + item[len(finalPath)+1:] + "\n")
		}
		stringBuffer.WriteString("\n")
	}

	stringBuffer.WriteString("## 任务分解\n\n")
	stringBuffer.WriteString("根据要求继续完成设计工作\n")

	result = stringBuffer.String()
	return
}
