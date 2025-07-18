package biz

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AttemptCompletionTool struct {
	commander    Commander
	workspacePtr *models.Workspace
}

func (s *AttemptCompletionTool) Name() string { return "attempt_completion" }
func (s *AttemptCompletionTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("Failed to read prompts: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AttemptCompletionTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *AttemptCompletionTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	fullPath := s.workspacePtr.GetWorkingPath()
	if !fp.Exist(fullPath) {
		err = fmt.Errorf("workspace not exist: %s", fullPath)
		return
	}

	resultVal, resultOK := param["result"].(string)
	if !resultOK || resultVal == "" {
		err = fmt.Errorf("missing or invalid result parameter")
		return
	}
	command, _ := param["command"].(string)
	log.Infof("command: [%s]", command)
	command = strings.TrimSpace(command)
	if command != "" {
		runResult, runErr := s.commander.RunCommand(s.workspacePtr.Name, fullPath, command)
		if runErr != nil {
			log.Errorf("run command failed, error:%s", runErr.Error())
			err = errors.New(runErr.Error())
			return
		}
		result = runResult
		return
	}

	result = "attempt completion successed."
	return
}
