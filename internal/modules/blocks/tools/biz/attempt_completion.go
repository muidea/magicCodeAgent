package biz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AttemptCompletionTool struct {
	commander Commander
}

func (s *AttemptCompletionTool) Name() string { return "attempt_completion" }
func (s *AttemptCompletionTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("Failed to read prompts: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AttemptCompletionTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *AttemptCompletionTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
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
			ret.Final = true
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	resultVal, resultOK := param["result"].(string)
	if !resultOK || resultVal == "" {
		err = fmt.Errorf("missing or invalid result parameter")
		return
	}
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	command, _ := param["command"].(string)
	log.Infof("command: [%s]", command)
	if command != "" {
		var args []string
		items := strings.Split(strings.TrimSpace(command), " ")
		if len(items) > 0 {
			command = items[0]
			args = items[1:]
			runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, args)
			if runErr != nil {
				err = errors.New(runErr.Error())
				return
			}
			result = runResult
			return
		}
	}
	result = resultVal
	return
}
