package biz

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
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

func (s *DecomposeDesignSubtaskTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	siblingsList, childList, workPath := s.getSiblingsAndChilds(workspaceDir, stage, name, level)

	planFileName := fmt.Sprintf("decompose-design-plan-%s", level)
	planContent := s.promptsRegistry.ReadSectionPromptContent(planFileName)
	planParam := map[string]string{
		"Name":             name,
		"Description":      description,
		"Siblings":         siblingsList,
		"Childrens":        childList,
		"WorkingDirectory": workPath,
	}

	tmpl, err := template.New(planFileName).Parse(planContent)
	if err != nil {
		log.Errorf("parse decompose-design-plan template failed:%s", err.Error())
		return
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, planParam)
	if err != nil {
		log.Errorf("execute decompose-design-plan template failed:%s", err.Error())
		return
	}

	result = builder.String()
	return
}

func (s *DecomposeDesignSubtaskTool) getSiblingsAndChilds(workspaceDir, stage, name, level string) (siblings, childs, workPath string) {
	rootPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(rootPath) {
		os.MkdirAll(rootPath, 0755)
	}
	fileList, fileErr := fp.ListPath(rootPath, "*.md", true)
	if fileErr != nil {
		return
	}

	log.Infof("rootPath:\n%s\n%s\n", rootPath, fileList)

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
	case "point":
		offset = 3
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

	siblingBuffer := strings.Builder{}
	for _, item := range siblingList {
		siblingBuffer.WriteString("- " + item + "\n")
	}
	siblingBuffer.WriteString("\n")
	siblings = siblingBuffer.String()

	childBuffer := strings.Builder{}
	for _, item := range childList {
		childBuffer.WriteString("- " + item + "\n")
	}
	childBuffer.WriteString("\n")
	childs = childBuffer.String()

	workPath = currentPath

	log.Infof("slibling:\n%s \nchild:\n%s\nworkPath:%s\n", siblings, childs, workPath)
	return
}
