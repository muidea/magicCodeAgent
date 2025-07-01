package biz

import (
	"fmt"
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

type DecomposeDesignPlanTool struct {
	promptsRegistry iipc.PromptsRegistry
}

func (s *DecomposeDesignPlanTool) Name() string {
	return "decompose_design_plan"
}

func (s *DecomposeDesignPlanTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeDesignPlanTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeDesignPlanTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Name: s.Name(),
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

	parentPath := filepath.Dir(workspaceDir)
	siblingsList, childList, workPath := s.getSiblingsAndChilds(parentPath, stage, name, level)

	planParam := map[string]string{
		"Name":             name,
		"Description":      description,
		"Path":             pathVal,
		"Siblings":         siblingsList,
		"Childrens":        childList,
		"WorkingDirectory": workPath,
	}
	planFileName := fmt.Sprintf("decompose-design-plan-%s-%s", level, workspaceType)
	result = s.promptsRegistry.ReadSectionPromptContent(planFileName, planParam)
	return
}

func (s *DecomposeDesignPlanTool) getSiblingsAndChilds(workspaceDir, stage, name, level string) (siblings, childs, workPath string) {
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
