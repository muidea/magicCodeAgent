package biz

import (
	"encoding/json"
	"fmt"
	"os"
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

type InsertContentTool struct{}

func (s *InsertContentTool) Name() string { return "insert_content" }
func (s *InsertContentTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *InsertContentTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *InsertContentTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid file path parameter")
		return
	}

	operations, ok := param["operations"].(string)
	if !ok || operations == "" {
		err = fmt.Errorf("missing or invalid operation parameter")
		return
	}

	var ops []struct {
		StartLine int    `json:"start_line"`
		Content   string `json:"content"`
	}

	if err = json.Unmarshal([]byte(operations), &ops); err != nil {
		err = fmt.Errorf("failed to parse operations: %v", err)
		return
	}

	if len(ops) == 0 {
		err = fmt.Errorf("operations cannot be empty")
		return
	}
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	var insertBlocks []*InsertBlock
	for _, op := range ops {
		insertBlocks = append(insertBlocks, &InsertBlock{
			StartLine: op.StartLine,
			Content:   op.Content,
		})
	}

	if err = InsertContent(fullPath, path, insertBlocks); err != nil {
		log.Errorf("InsertContent failed, fullPath:%s, path:%s err: %v", fullPath, path, err)
		return
	}

	result = "insert content successfully!"
	return
}

/*
## insert_content
Description: Inserts content at specific line positions in a file. This is the primary tool for adding new content and code (functions/methods/classes, imports, attributes etc.) as it allows for precise insertions without overwriting existing content. The tool uses an efficient line-based insertion system that maintains file integrity and proper ordering of multiple insertions. Beware to use the proper indentation. This tool is the preferred way to add new content and code to files.
描述：在文件的特定行位置插入内容。这是添加新内容和代码（函数/方法/类/导入/属性等）的主要工具，因为它允许精确插入而不会覆盖现有内容。该工具使用高效的基于行的插入系统，保持文件完整性和多个插入的正确顺序。注意使用正确的缩进。这是向文件添加新内容和代码的首选方式。

Parameters 参数:
- path: (required) The path of the file to insert content into (relative to the current workspace directory /home/rangh/webspace)
- path: (必填) 要插入内容的文件路径（相对于当前工作空间目录 /home/rangh/webspace）

- operations: (required) A JSON array of insertion operations. Each operation is an object with:
- operations: (必填) 插入操作的JSON数组。每个操作都是一个对象，包含：
    * start_line: (required) The line number where the content should be inserted. The content currently at that line will end up below the inserted content.
    * start_line: (必填) 内容应该插入的行号。该行当前的内容将会出现在插入内容的下方。
    * content: (required) The content to insert at the specified position. IMPORTANT NOTE: If the content is a single line, it can be a string. If it's a multi-line content, it should be a string with newline characters (
) for line breaks. Make sure to include the correct indentation for the content.
    * content: (必填) 要在指定位置插入的内容。重要提示：如果内容是单行的，可以是字符串。如果是多行内容，应该是包含换行符的字符串。确保包含正确的内容缩进。

Usage 使用方法:
<insert_content>
<path>File path here</path>
<operations>[
  {
    "start_line":10,
    "content":"Your content here"
  }
]</operations>
</insert_content>

Example 示例: Insert a new function and its import statement 插入新函数及其导入语句
<insert_content>
<path>File path here</path>
<operations>[
  {
    "start_line":1,
    "content":"import { sum } from './utils';"
  },
  {
    "start_line":10,
    "content":"function calculateTotal(items: number[]): number {
    return items.reduce((sum, item) => sum + item, 0);
}"
  }
]</operations>
</insert_content>

*/

type InsertBlock struct {
	StartLine int    `json:"start_line"`
	Content   string `json:"content"`
}

// InsertContent 实现insert_content工具功能
func InsertContent(workspaceDir, path string, operations []*InsertBlock) error {
	fullPath := filepath.Join(workspaceDir, path)

	// 读取原始文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// 按start_line降序排序操作，以便从后往前插入
	for i := len(operations) - 1; i >= 0; i-- {
		op := operations[i]
		if op.StartLine < 1 || op.StartLine > len(lines)+1 {
			return fmt.Errorf("invalid start_line: %d", op.StartLine)
		}

		insertPos := op.StartLine - 1
		insertLines := strings.Split(op.Content, "\n")

		// 插入新内容
		newLines := make([]string, 0, len(lines)+len(insertLines))
		newLines = append(newLines, lines[:insertPos]...)
		newLines = append(newLines, insertLines...)
		newLines = append(newLines, lines[insertPos:]...)
		lines = newLines
	}

	// 写入文件
	return os.WriteFile(fullPath, []byte(strings.Join(lines, "\n")), 0644)
}
