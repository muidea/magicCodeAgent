package biz

import (
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

type ExecuteCommandTool struct {
	commander    Commander
	workspacePtr *models.Workspace
}

func (s *ExecuteCommandTool) Name() string { return "execute_command" }
func (s *ExecuteCommandTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ExecuteCommandTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ExecuteCommandTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	command, ok := param["command"].(string)
	if !ok || command == "" {
		err = fmt.Errorf("command is empty")
		return
	}
	cwd, _ := param["cwd"].(string)

	fullPath = filepath.Join(fullPath, cwd)
	var args []string
	if argsParam, ok := param["args"].([]interface{}); ok {
		for _, arg := range argsParam {
			if strArg, ok := arg.(string); ok {
				args = append(args, strArg)
			}
		}
	}
	command = fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	runResult, runErr := s.commander.RunCommand(s.workspacePtr.Name, fullPath, command)
	if runErr != nil {
		log.Errorf("run command failed, error: %s", runErr.Error())
		result = fmt.Sprintf("execute command failed, error:%s", runErr.Error())
		return
	}
	result = "execute command finished."
	if runResult != "" {
		result = fmt.Sprintf("%s\n%s", result, runResult)
	}
	return
}
