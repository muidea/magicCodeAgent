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

type UpdateTodoListTool struct{}

func (s *UpdateTodoListTool) Name() string { return "update_todo_list" }
func (s *UpdateTodoListTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *UpdateTodoListTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *UpdateTodoListTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	todos, ok := param["todos"].(string)
	if !ok || todos == "" {
		err = fmt.Errorf("missing or invalid todos parameter")
		return
	}

	result = "todo list updated successfully"
	return
}
