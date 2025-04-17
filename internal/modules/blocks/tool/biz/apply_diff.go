package biz

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
## apply_diff
Description: Request to replace existing code using a search and replace block.
This tool allows for precise, surgical replaces to files by specifying exactly what content to search for and what to replace it with.
The tool will maintain proper indentation and formatting while making changes.
Only a single operation is allowed per tool use.
The SEARCH section must exactly match existing content including whitespace and indentation.
If you're not confident in the exact content to search for, use the read_file tool first to get the exact content.
When applying the diffs, be extra careful to remember to change any closing brackets or other syntax that may be affected by the diff farther down in the file.
ALWAYS make as many changes in a single 'apply_diff' request as possible using multiple SEARCH/REPLACE blocks

描述：使用搜索替换块请求替换现有代码。该工具通过精确指定要搜索的内容和替换内容，实现文件的精准修改。工具在更改时会保持正确的缩进和格式。每次工具使用仅允许单个操作。SEARCH部分必须完全匹配现有内容（包括空格和缩进）。如果不确定要搜索的确切内容，请先使用read_file工具获取准确内容。应用差异时要特别注意修改可能受影响的闭合括号或其他语法。尽量在单个apply_diff请求中使用多个SEARCH/REPLACE块完成所有更改。

Parameters:
- path: (required) The path of the file to modify (relative to the current workspace directory /home/rangh/webspace)
  路径：(必填) 要修改的文件路径（相对于当前工作空间目录/home/rangh/webspace）
- diff: (required) The search/replace block defining the changes.
  差异内容：(必填) 定义更改的搜索/替换块

Diff format:
```
<<<<<<< SEARCH
:start_line: (required) The line number of original content where the search block starts.
:end_line: (required) The line number of original content  where the search block ends.
-------
[exact content to find including whitespace]
=======
[new content to replace with]
>>>>>>> REPLACE

```


Example:
示例：

Original file:
```
1 | def calculate_total(items):
2 |     total = 0
3 |     for item in items:
4 |         total += item
5 |     return total
```

Search/Replace content:
```
<<<<<<< SEARCH
:start_line:1
:end_line:5
-------
def calculate_total(items):
    total = 0
    for item in items:
        total += item
    return total
=======
def calculate_total(items):
    \"\"\"Calculate total with 10% markup\"\"\"
    return sum(item * 1.1 for item in items)
>>>>>>> REPLACE

```

Search/Replace content with multi edits:
```
<<<<<<< SEARCH
:start_line:1
:end_line:2
-------
def calculate_total(items):
    sum = 0
=======
def calculate_sum(items):
    sum = 0
>>>>>>> REPLACE

<<<<<<< SEARCH
:start_line:4
:end_line:5
-------
        total += item
    return total
=======
        sum += item
    return sum
>>>>>>> REPLACE
```


Usage:
<apply_diff>
<path>File path here</path>
<diff>
Your search/replace content here
You can use multi search/replace block in one diff block, but make sure to include the line numbers for each block.
Only use a single line of '=======' between search and replacement content, because multiple '=======' will corrupt the file.
</diff>
</apply_diff>

*/
// ApplyDiff 实现apply_diff工具功能
func (s *Tool) ApplyDiff(path string, diff string) error {
	fullPath := filepath.Join(WorkspaceDir, path)

	// 读取原始文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	// 解析diff内容
	diffParts := strings.Split(diff, "=======")
	if len(diffParts) != 2 {
		return fmt.Errorf("invalid diff format")
	}

	searchContent := strings.TrimSpace(diffParts[0])
	replaceContent := strings.TrimSpace(diffParts[1])

	// 替换内容
	newContent := bytes.Replace(content, []byte(searchContent), []byte(replaceContent), 1)

	// 写入文件
	return os.WriteFile(fullPath, newContent, 0644)
}
