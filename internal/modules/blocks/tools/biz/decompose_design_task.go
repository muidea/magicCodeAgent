package biz

import (
	"fmt"
	"html/template"
	"os"
	"path"
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

type DecomposeDesignTaskTool struct {
	workspacePtr *models.Workspace
	templatePtr  *models.Template
}

func (s *DecomposeDesignTaskTool) Name() string {
	return "decompose_design_task"
}

func (s *DecomposeDesignTaskTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeDesignTaskTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeDesignTaskTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Name: s.Name(),
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
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
	pathVal := ""
	path, ok := param["path"].(string)
	if ok {
		pathVal = path
	}
	description, ok := param["description"].(string)
	if !ok || description == "" {
		err = fmt.Errorf("description is empty")
		return
	}

	rootPath := s.workspacePtr.GetRootPath()
	stage := "design"
	siblingsList, childList, workPath := s.getSiblingsAndChilds(rootPath, stage, name, level)

	planParam := map[string]string{
		"Name":             name,
		"Description":      description,
		"ReferencePath":    pathVal,
		"Siblings":         siblingsList,
		"Childrens":        childList,
		"WorkingDirectory": workPath,
	}
	planFileName := fmt.Sprintf("decompose-design-task-%s", level)
	promptVal, promptErr := s.templatePtr.ReadPrompt(planFileName)
	if promptErr != nil {
		err = promptErr
		return
	}

	tmplPtr, tmplErr := template.New(planFileName).Parse(promptVal)
	if tmplErr != nil {
		err = tmplErr
		log.Errorf("parse system info template failed:%s", tmplErr.Error())
		return
	}

	var builder strings.Builder
	err = tmplPtr.Execute(&builder, planParam)
	if err != nil {
		log.Errorf("execute system info template failed:%s", err.Error())
		return
	}

	result = builder.String()
	return
}

func (s *DecomposeDesignTaskTool) getSiblingsAndChilds(workspaceDir, stage, name, level string) (siblings, childs, workPath string) {
	rootPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(rootPath) {
		os.MkdirAll(rootPath, 0755)
	}
	fileList, fileErr := fp.ListPath(rootPath, "*.md", true)
	if fileErr != nil {
		return
	}

	currentPath := "."
	childPath := "."
	siblingList := fu.StringSet{}
	childList := fu.StringSet{}

	offset := -1
	switch level {
	case models.ProjectLevel:
		offset = 0
	case models.ModuleLevel:
		offset = 1
	case models.BlockLevel:
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
		currentPath = path.Join(rootPath, filepath.Dir(childPath))
		siblingList, _ = fp.ListPath(currentPath, "*.md", false)
	}

	siblingBuffer := strings.Builder{}
	for _, item := range siblingList {
		item := item[len(rootPath)+1:]
		siblingBuffer.WriteString("- " + item + "\n")
	}
	siblings = siblingBuffer.String()

	childBuffer := strings.Builder{}
	for _, item := range childList {
		childBuffer.WriteString("- " + item + "\n")
	}
	childs = childBuffer.String()

	workPath = currentPath

	//log.Infof("slibling:\n%s \nchild:\n%s\nworkPath:%s\n", siblings, childs, workPath)
	return
}
