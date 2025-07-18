package biz

import (
	"fmt"
	"html/template"
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

type DecomposeCodePlanTool struct {
	workspacePtr *models.Workspace
	templatePtr  *models.Template
}

func (s *DecomposeCodePlanTool) Name() string {
	return "decompose_code_plan"
}

func (s *DecomposeCodePlanTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeCodePlanTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeCodePlanTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	childList, childErr := s.getChildFiles(rootPath, "design", pathVal)
	if childErr != nil {
		log.Errorf("failed to get child files, error:%s", childErr.Error())
		result = "error: 制定开发目标失败!\n" + childErr.Error()
		return
	}

	planParam := map[string]string{
		"DesignPath": childList,
	}
	planFileName := fmt.Sprintf("decompose-code-plan-%s", level)
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

func (s *DecomposeCodePlanTool) getChildFiles(workspaceDir, stage, pathVal string) (childsFiles string, err error) {
	rootPath := filepath.Join(workspaceDir, stage)
	finalParentFile := filepath.Join(rootPath, pathVal)
	if !fp.Exist(finalParentFile) {
		log.Errorf("finalParentFile:%s not exist", finalParentFile)
		err = fmt.Errorf("设计文档:%s不存在, 需要重新核对路径", pathVal)
		return
	}

	childList := fu.StringSet{}
	filePath := finalParentFile[len(rootPath)+1:]
	childList = childList.Add(filePath)

	parentFilePath, parentFileName := path.Split(finalParentFile)
	subItems := strings.Split(parentFileName, "-")
	if len(subItems) != 2 {
		err = fmt.Errorf("设计文档:%s不存在, 需要重新核对路径", pathVal)
		return
	}

	finalChildPath := filepath.Join(parentFilePath, subItems[0])
	if fp.Exist(finalChildPath) {
		childFileList, childFileErr := fp.ListPath(finalChildPath, "*.md", false)
		if childFileErr != nil {
			log.Errorf("list path:%s error:%s", finalChildPath, childFileErr)
			err = fmt.Errorf("访问设计文档: %s 对应的子目录出错: %s", pathVal, childFileErr.Error())
			return
		}

		for _, val := range childFileList {
			filePath = val[len(rootPath)+1:]
			childList = childList.Add(filePath)
		}
	}

	childBuffer := strings.Builder{}
	for _, item := range childList {
		childBuffer.WriteString("- " + item + "\n")
	}
	childsFiles = childBuffer.String()
	log.Infof("child:\n%s\n", childsFiles)
	return
}
