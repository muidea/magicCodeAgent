package biz

import (
	"fmt"
	"os"
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

type DecomposeDesignSubtaskTool struct {
	promptsRegistry iipc.PromptsRegistry
}

func (s *DecomposeDesignSubtaskTool) Name() string {
	return "decompose_design_subtask"
}

func (s *DecomposeDesignSubtaskTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeDesignSubtaskTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeDesignSubtaskTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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
	stringBuffer.WriteString("## " + name + " - 分解计划\n\n")
	stringBuffer.WriteString("### 待分解功能描述\n\n")
	stringBuffer.WriteString(description + "\n\n")

	// ### 任务分解要求
	levelSectionName := fmt.Sprintf("decompose-design-%s-plan", level)
	content := s.promptsRegistry.ReadSectionPromptContent(levelSectionName)
	if content != "" {
		stringBuffer.WriteString(content)
		stringBuffer.WriteString("\n\n")
	}

	// 查看兄弟层级分解清单
	rootPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(rootPath) {
		os.MkdirAll(rootPath, 0755)
	}
	fileList, fileErr := fp.ListPath(rootPath, "*.md", true)
	if fileErr != nil {
		err = fileErr
		return
	}

	currentPath := "."
	childPath := "."
	siblingList := fu.StringSet{}
	childList := fu.StringSet{}

	offset := -1
	switch level {
	case "project":
		offset = 0
	case "module":
		offset = 1
	case "block":
		offset = 2
	}

	for _, val := range fileList {
		filePath := val[len(rootPath)+1:]
		items := strings.Split(filePath, "/")
		if len(items) != offset+1 {
			continue
		}
		if offset == 0 {
			childList = childList.Add(filePath)
			continue
		}
		if items[offset-1] != name {
			continue
		}
		childPath = strings.Join(items[:offset], "/")
		childList = childList.Add(filePath)
	}
	if childPath != "." {
		currentPath = filepath.Dir(childPath)
		siblingList, _ = fp.ListPath(currentPath, "*.md", false)
	}

	if level != "project" {
		stringBuffer.WriteString("### 同级功能清单\n\n")
		if len(siblingList) > 0 {
			for _, item := range siblingList {
				stringBuffer.WriteString("- " + item + "\n")
			}
			stringBuffer.WriteString("\n")
		}
	}
	stringBuffer.WriteString("### 已分解子功能清单\n\n")
	if len(fileList) > 0 {
		for _, item := range childList {
			stringBuffer.WriteString("- " + item + "\n")
		}
		stringBuffer.WriteString("\n")
	}

	stringBuffer.WriteString("### 工作目录信息\n\n")
	stringBuffer.WriteString("当前分解任务所在目录：" + currentPath + "\n")
	stringBuffer.WriteString("子功能任务所在目录：" + childPath + "\n")

	result = stringBuffer.String()
	return
}
