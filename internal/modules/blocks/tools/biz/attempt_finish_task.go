package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AttemptFinishTaskTool struct{}

func (s *AttemptFinishTaskTool) Name() string { return "attempt_finish_task" }
func (s *AttemptFinishTaskTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("Failed to read prompts: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AttemptFinishTaskTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *AttemptFinishTaskTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("path is empty")
		return
	}
	resultVal, resultOK := param["result"].(string)
	if !resultOK || resultVal == "" {
		err = fmt.Errorf("missing or invalid result parameter")
		return
	}

	result = "attempt finish task successed."
	return
}
