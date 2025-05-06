package biz

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/task"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/config"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

type Tools struct {
	biz.Base

	toolRegistry *ToolRegistry
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
	promptsRegistry iipc.PromptsRegistry,
) *Tools {
	ptr := &Tools{
		Base:         biz.New(imbtc.ToolsModule, eventHub, backgroundRoutine),
		toolRegistry: NewToolRegistry(promptsRegistry),
	}

	ptr.toolRegistry.Register(&ReadFileTool{})
	ptr.toolRegistry.Register(&FetchInstructionsTool{})
	ptr.toolRegistry.Register(&SearchFilesTool{})
	ptr.toolRegistry.Register(&ListFilesTool{})
	ptr.toolRegistry.Register(&ListCodeDefinitionNamesTool{})
	ptr.toolRegistry.Register(&ApplyDiffTool{})
	ptr.toolRegistry.Register(&WriteToFileTool{})
	ptr.toolRegistry.Register(&InsertContentTool{})
	ptr.toolRegistry.Register(&SearchAndReplaceTool{})
	ptr.toolRegistry.Register(&ExecuteCommandTool{})
	ptr.toolRegistry.Register(&UseMcpToolTool{})
	ptr.toolRegistry.Register(&AccessMcpResourceTool{})
	ptr.toolRegistry.Register(&AskFollowupQuestionTool{})
	ptr.toolRegistry.Register(&AttemptCompletionTool{})
	ptr.toolRegistry.Register(&SwitchModeTool{})
	ptr.toolRegistry.Register(&NewTaskTool{})

	ptr.SubscribeFunc(imbtc.ExecuteTool, ptr.Execute)
	return ptr
}

func (s *Tools) Run() (err *cd.Error) {
	currentWorkDir := fu.GetCurrentWorkDir()

	dataFile := config.GetDataFile()
	if dataFile != "" {
		s.execDataFile(currentWorkDir, dataFile)
		return
	}

	dataContent := config.GetDataContent()
	if dataContent != "" {
		s.execDataContent(currentWorkDir, dataContent)
		return
	}

	return
}

func (s *Tools) execDataFile(workspaceDir, dataFile string) (result imbtc.ToolResult, err *cd.Error) {
	if !fp.Exist(dataFile) {
		err = cd.NewError(cd.IllegalParam, "data file not exist")
		log.Errorf("data file not exist: %s", dataFile)
		return
	}

	byteContent, byteErr := os.ReadFile(dataFile)
	if byteErr != nil {
		err = cd.NewError(cd.IllegalParam, "read data file error")
		log.Errorf("read data file error: %s", byteErr.Error())
		return
	}

	result, err = s.execDataContent(workspaceDir, string(byteContent))
	return
}

func (s *Tools) execDataContent(workspaceDir, dataContent string) (result imbtc.ToolResult, err *cd.Error) {
	dataBlocks, dataErr := util.SpliteBlocks(dataContent)
	if dataErr != nil {
		err = cd.NewError(cd.IllegalParam, "split data content error")
		log.Errorf("split data content error: %s", dataErr.Error())
		return
	}

	byteBuffer := bytes.NewBuffer(nil)
	for idx, val := range dataBlocks {
		toolName, toolParam, toolErr := parseParamContent(val)
		if toolName == "thinking" {
			continue
		}

		if toolErr != nil {
			err = cd.NewError(cd.IllegalParam, "parse param content error")
			log.Errorf("parse param content error: %s", toolErr.Error())
			return
		}
		if toolName == "" {
			err = cd.NewError(cd.IllegalParam, "tool name is empty")
			log.Errorf("tool name is empty")
			return
		}

		toolResult, toolErr := s.toolRegistry.Execute(workspaceDir, toolName, toolParam)
		if toolErr != nil {
			err = cd.NewError(cd.IllegalParam, "execute tool error")
			log.Errorf("execute tool error: %s", toolErr.Error())
			return
		}
		if idx > 0 {
			byteBuffer.WriteString("======\n\n")
		}
		byteBuffer.Write(toolResult)
	}

	result = byteBuffer.Bytes()
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
		err = fmt.Errorf("invalid tool content")
		return
	}

	return
}

func (s *Tools) Execute(ev event.Event, re event.Result) {
	var err *cd.Error
	var executeResult imbtc.ToolResult
	defer func() {
		if re != nil {
			re.Set(executeResult, err)
		}
	}()

	dataVal := ev.Data()
	contentVal, contentOK := dataVal.(string)
	if !contentOK {
		err = cd.NewError(cd.IllegalParam, "execute tool failed, invalid tool content")
		log.Warnf("content is not string, dataVal:%v", dataVal)
		return
	}
	workspaceDirVal := ev.GetData("workspaceDir")
	workspaceDir, workspaceDirOK := workspaceDirVal.(string)
	if !workspaceDirOK {
		err = cd.NewError(cd.IllegalParam, "execute tool failed, workspaceDir is not string")
		log.Warnf("workspaceDir is not string, dataVal:%v", dataVal)
		return
	}

	executeResult, err = s.execDataContent(workspaceDir, contentVal)
}
