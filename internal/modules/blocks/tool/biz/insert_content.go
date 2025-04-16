package biz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
## insert_content
Description: Inserts content at specific line positions in a file. This is the primary tool for adding new content and code (functions/methods/classes, imports, attributes etc.) as it allows for precise insertions without overwriting existing content. The tool uses an efficient line-based insertion system that maintains file integrity and proper ordering of multiple insertions. Beware to use the proper indentation. This tool is the preferred way to add new content and code to files.
Parameters:
- path: (required) The path of the file to insert content into (relative to the current workspace directory /home/rangh/webspace)
- operations: (required) A JSON array of insertion operations. Each operation is an object with:
    * start_line: (required) The line number where the content should be inserted.  The content currently at that line will end up below the inserted content.
    * content: (required) The content to insert at the specified position. IMPORTANT NOTE: If the content is a single line, it can be a string. If it's a multi-line content, it should be a string with newline characters (
) for line breaks. Make sure to include the correct indentation for the content.
Usage:
<insert_content>
<path>File path here</path>
<operations>[
  {
    \"start_line\":10,
    \"content\":\"Your content here\"
  }
]</operations>
</insert_content>
Example: Insert a new function and its import statement
<insert_content>
<path>File path here</path>
<operations>[
  {
    \"start_line\":1,
    \"content\":\"import { sum } from './utils';\"
  },
  {
    \"start_line\":10,
    \"content\":\"function calculateTotal(items: number[]): number {
    return items.reduce((sum, item) => sum + item, 0);
}\"
  }
]</operations>
</insert_content>

*/

// InsertContent 实现insert_content工具功能
func (s *Tool) InsertContent(path string, operations []struct {
	StartLine int    `json:"start_line"`
	Content   string `json:"content"`
}) error {
	fullPath := filepath.Join(WorkspaceDir, path)
	
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