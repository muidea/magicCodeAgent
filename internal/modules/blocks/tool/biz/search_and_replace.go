package biz

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
## search_and_replace
Description: Request to perform search and replace operations on a file. Each operation can specify a search pattern (string or regex) and replacement text, with optional line range restrictions and regex flags. Shows a diff preview before applying changes.
Parameters:
- path: (required) The path of the file to modify (relative to the current workspace directory /home/rangh/webspace)
- operations: (required) A JSON array of search/replace operations. Each operation is an object with:
    * search: (required) The text or pattern to search for
    * replace: (required) The text to replace matches with. If multiple lines need to be replaced, use \"
\" for newlines
    * start_line: (optional) Starting line number for restricted replacement
    * end_line: (optional) Ending line number for restricted replacement
    * use_regex: (optional) Whether to treat search as a regex pattern
    * ignore_case: (optional) Whether to ignore case when matching
    * regex_flags: (optional) Additional regex flags when use_regex is true
Usage:
<search_and_replace>
<path>File path here</path>
<operations>[
  {
    \"search\":\"text to find\",
    \"replace\":\"replacement text\",
    \"start_line\":1,
    \"end_line\":10
  }
]</operations>
</search_and_replace>
Example: Replace \"foo\" with \"bar\" in lines 1-10 of example.ts
<search_and_replace>
<path>example.ts</path>
<operations>[
  {
    \"search\":\"foo\",
    \"replace\":\"bar\",
    \"start_line\":1,
    \"end_line\":10
  }
]</operations>
</search_and_replace>
Example: Replace all occurrences of \"old\" with \"new\" using regex
<search_and_replace>
<path>example.ts</path>
<operations>[
  {
    \"search\":\"old\\w+\",
    \"replace\":\"new$&\",
    \"use_regex\":true,
    \"ignore_case\":true
  }
]</operations>
</search_and_replace>

*/

// SearchAndReplace 实现search_and_replace工具功能
func (s *Tool) SearchAndReplace(path string, operations []struct {
	Search      string `json:"search"`
	Replace     string `json:"replace"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	UseRegex    bool   `json:"use_regex"`
	IgnoreCase  bool   `json:"ignore_case"`
	RegexFlags  string `json:"regex_flags"`
}) error {
	fullPath := filepath.Join(WorkspaceDir, path)
	
	// 读取原始文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	
	lines := strings.Split(string(content), "\n")
	
	for _, op := range operations {
		start := op.StartLine - 1
		if start < 0 {
			start = 0
		}
		end := op.EndLine
		if end <= 0 || end > len(lines) {
			end = len(lines)
		}
		
		for i := start; i < end; i++ {
			var newLine string
			if op.UseRegex {
				regex := regexp.MustCompile(op.Search)
				newLine = regex.ReplaceAllString(lines[i], op.Replace)
			} else {
				newLine = strings.ReplaceAll(lines[i], op.Search, op.Replace)
			}
			lines[i] = newLine
		}
	}
	
	// 写入文件
	return os.WriteFile(fullPath, []byte(strings.Join(lines, "\n")), 0644)
}