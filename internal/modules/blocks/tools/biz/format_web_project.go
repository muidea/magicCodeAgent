package biz

import (
	"fmt"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type FormatWebProjectTool struct {
	commander Commander
}

func (s *FormatWebProjectTool) Name() string {
	return "format_web_project"
}

func (s *FormatWebProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *FormatWebProjectTool) Tags() fu.StringSet { return fu.StringSet{iipc.Web} }

func (s *FormatWebProjectTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	fullPath := workspaceDir
	if !fp.Exist(fullPath) {
		err = fmt.Errorf("workspace not exist: %s", fullPath)
		return
	}

	// yarn run format  --write --ignore-path .prettierignore
	command := "yarn"
	var args = []string{
		"run", "format",
		"--write",
		"--ignore-path", ".prettierignore",
	}
	runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, 0, args)
	if runErr != nil {
		log.Errorf("run command failed, error:%s", runErr.Error())
		result = fmt.Sprintf("format web project failed, error:%s", runErr)
		return
	}
	result = "format web project finished."
	if runResult != "" {
		result = fmt.Sprintf("%s\n%s", result, runResult)
	}
	return
}
