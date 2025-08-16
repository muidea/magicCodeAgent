package biz

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type DecomposeDebugTaskTool struct {
	workspacePtr *models.Workspace
	templatePtr  *models.Template
}

func (s *DecomposeDebugTaskTool) Name() string {
	return "decompose_debug_task"
}

func (s *DecomposeDebugTaskTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DecomposeDebugTaskTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *DecomposeDebugTaskTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	//path, ok := param["path"].(string)
	//if ok {
	//	pathVal = path
	//}
	description, ok := param["description"].(string)
	if !ok || description == "" {
		err = fmt.Errorf("description is empty")
		return
	}

	planParam := map[string]string{
		"Name":        name,
		"Description": description,
	}
	planFileName := fmt.Sprintf("decompose-debug-task-%s", "project")
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
