package biz

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

var WorkspaceDir = "/home/rangh/codespace/magicCodeAgent/internal/modules/blocks/tool/biz/temp"

func init() {
	WorkspaceDir, _ = os.Getwd()

	WorkspaceDir = path.Join(WorkspaceDir, "temp")
}

type Tools struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
	promptsRegistry iipc.PromptsRegistry,
) *Tools {
	ptr := &Tools{
		Base: biz.New(common.ToolsModule, eventHub, backgroundRoutine),
	}

	if promptsRegistry != nil {
		promptsRegistry.AddPrompts("tools", &ReadFileTool{})
		promptsRegistry.AddPrompts("tools", &FetchInstructionsTool{})
		promptsRegistry.AddPrompts("tools", &SearchFilesTool{})
		promptsRegistry.AddPrompts("tools", &ListFilesTool{})
		promptsRegistry.AddPrompts("tools", &ListCodeDefinitionNamesTool{})
		promptsRegistry.AddPrompts("tools", &ApplyDiffTool{})
		promptsRegistry.AddPrompts("tools", &WriteToFileTool{})
		promptsRegistry.AddPrompts("tools", &InsertContentTool{})
		promptsRegistry.AddPrompts("tools", &SearchAndReplaceTool{})
		promptsRegistry.AddPrompts("tools", &ExecuteCommandTool{})
		promptsRegistry.AddPrompts("tools", &UseMcpToolTool{})
		promptsRegistry.AddPrompts("tools", &AccessMcpResourceTool{})
		promptsRegistry.AddPrompts("tools", &AskFollowupQuestionTool{})
		promptsRegistry.AddPrompts("tools", &AttemptCompletionTool{})
		promptsRegistry.AddPrompts("tools", &SwitchModeTool{})
		promptsRegistry.AddPrompts("tools", &NewTaskTool{})
	}

	return ptr
}

func (s *Tools) Run() (err *cd.Error) {
	return
}

// 解析执行参数
// content里解析出执行函数名和参数
// content的格式是xml
// 执行函数是xml的root节点名
// 执行函数的参数是root的子节点及子节点的值
// 参照示例参见各个工具函数对应的注释信息
func parseParamContent(content string) (ret string, param map[string]any, err error) {
	// 解析XML内容
	decoder := xml.NewDecoder(strings.NewReader(content))
	param = make(map[string]any)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if ret == "" {
				ret = se.Name.Local
			} else {
				var value string
				if err := decoder.DecodeElement(&value, &se); err != nil {
					return "", nil, err
				}
				param[se.Name.Local] = value
			}
		}
	}

	if ret == "" {
		return "", nil, fmt.Errorf("invalid xml content")
	}

	return ret, param, nil
}
