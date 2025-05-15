package biz

import (
	"strings"
)

// findContentInFile 在文件内容中查找匹配的内容，返回匹配内容的起始行号（1-based）
// 如果找不到匹配的内容，返回-1
// 支持精确匹配和模糊匹配（忽略空白字符的差异）
func findContentInFile(fileLines []string, content string) int {
	if content == "" {
		return -1
	}

	contentLines := strings.Split(content, "\n")
	if len(contentLines) == 0 {
		return -1
	}

	// 单行内容直接查找
	if len(contentLines) == 1 {
		// 精确匹配
		for i, line := range fileLines {
			if line == contentLines[0] {
				return i + 1 // 返回1-based行号
			}
		}

		// 尝试忽略空白字符的匹配
		contentTrimmed := strings.TrimSpace(contentLines[0])
		if contentTrimmed != "" {
			for i, line := range fileLines {
				if strings.TrimSpace(line) == contentTrimmed {
					return i + 1 // 返回1-based行号
				}
			}
		}
		return -1
	}

	// 多行内容需要连续匹配
	// 1. 先尝试精确匹配
	for i := 0; i <= len(fileLines)-len(contentLines); i++ {
		matched := true
		for j := 0; j < len(contentLines); j++ {
			if i+j >= len(fileLines) || fileLines[i+j] != contentLines[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 1 // 返回1-based行号
		}
	}

	// 2. 如果精确匹配失败，尝试忽略空白字符的匹配
	// 预处理内容行，去除空白
	contentTrimmed := make([]string, len(contentLines))
	for i, line := range contentLines {
		contentTrimmed[i] = strings.TrimSpace(line)
	}

	for i := 0; i <= len(fileLines)-len(contentLines); i++ {
		matched := true
		for j := 0; j < len(contentLines); j++ {
			if i+j >= len(fileLines) || strings.TrimSpace(fileLines[i+j]) != contentTrimmed[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 1 // 返回1-based行号
		}
	}

	return -1
}
