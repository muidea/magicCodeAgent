package biz

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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