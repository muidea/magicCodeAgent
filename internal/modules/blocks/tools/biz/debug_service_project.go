package biz

import (
	"errors"
	"fmt"
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type DebugServiceProjectTool struct {
	commander Commander
}

func (s *DebugServiceProjectTool) Name() string {
	return "debug_service_project"
}

func (s *DebugServiceProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DebugServiceProjectTool) Tags() fu.StringSet { return fu.StringSet{iipc.Service} }

func (s *DebugServiceProjectTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
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

	var restart bool
	restartVal, restartOK := param["restart"]
	if restartOK {
		restart = restartVal.(bool)
	}

	var runStatusMsg, runErrorMsg string
	var runErr *cd.Error

	fullPath := filepath.Join(workspaceDir, stage)
	runStatusMsg, runErrorMsg, runErr = s.commander.RunCommand(workspaceName, fullPath, "go", restart, []string{"test", "-v", "./..."})
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	if runErrorMsg != "" {
		result = runErrorMsg
		return
	}

	result = runStatusMsg
	//
	// TODO 执行调试命令
	// lint
	// test
	// build
	// 这其中任意一个命令执行出错，就返回失败
	return
}
