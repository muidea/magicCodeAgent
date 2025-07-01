package biz

import (
	"errors"
	"fmt"
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

type DecomposeCodePlanTool struct {
	promptsRegistry iipc.PromptsRegistry
}

func (s *DecomposeCodePlanTool) Name() string {
	return "decompose_code_plan"
}

func (s *DecomposeCodePlanTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeCodePlanTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeCodePlanTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	childList, childErr := s.getChildFiles(parentPath, "design", pathVal)
	if childErr != nil {
		result = "error: 制定开发目标失败!\n" + childErr.Error()
		return
	}

	planParam := map[string]string{
		"DesignPath": childList,
	}
	planFileName := fmt.Sprintf("decompose-code-plan-%s", level)
	result = s.promptsRegistry.ReadSectionPromptContent(planFileName, planParam)
	return
}

func (s *DecomposeCodePlanTool) getChildFiles(workspaceDir, stage, pathVal string) (childsFiles string, err error) {
	rootPath := filepath.Join(workspaceDir, stage)
	finalPath := filepath.Join(rootPath, pathVal)
	if !fp.Exist(finalPath) {
		err = fmt.Errorf("设计文档对应目录:%s不存在, 需要重新核对路径", pathVal)
		return
	}
	fileList, fileErr := fp.ListPath(finalPath, "*.md", false)
	if fileErr != nil {
		err = fmt.Errorf("list file path:%s error:%s", pathVal, fileErr.Error())
		return
	}

	childList := fu.StringSet{}
	for _, val := range fileList {
		filePath := val[len(rootPath)+1:]
		childList = childList.Add(filePath)
	}

	// 路径错误，找不不到设计文档
	if len(childList) == 0 {
		err = errors.New("not found design document in path:" + pathVal)
		return
	}

	childBuffer := strings.Builder{}
	for _, item := range childList {
		childBuffer.WriteString("- " + item + "\n")
	}
	childsFiles = childBuffer.String()
	log.Infof("child:\n%s\n", childsFiles)
	return
}
