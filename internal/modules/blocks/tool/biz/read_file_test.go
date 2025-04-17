package biz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	// 创建一个Tool实例用于测试
	eventHub := event.NewHub(100)
	backgroundRoutine := task.NewBackgroundRoutine(100)
	tool := New(eventHub, backgroundRoutine)

	// 保存原始的WorkspaceDir
	originalWorkspaceDir := WorkspaceDir
	// 临时修改WorkspaceDir为测试目录
	tempDir := t.TempDir()
	WorkspaceDir = tempDir
	// 确保测试结束后恢复原始值
	defer func() { WorkspaceDir = originalWorkspaceDir }()

	// 创建测试文件
	testFilePath := filepath.Join(tempDir, "test_file.txt")
	testContent := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")

	// 测试用例1：读取整个文件
	result1, err := tool.ReadFile("test_file.txt", 0, 0)
	assert.NoError(t, err, "读取整个文件应该成功")
	expectedLines := 5
	actualLines := strings.Count(result1, "\n")
	assert.Equal(t, expectedLines, actualLines, "应该读取到5行内容")
	assert.Contains(t, result1, "1 | Line 1", "输出应该包含行号和内容")
	assert.Contains(t, result1, "5 | Line 5", "输出应该包含最后一行")

	// 测试用例2：读取部分文件（指定开始行）
	result2, err := tool.ReadFile("test_file.txt", 3, 0)
	assert.NoError(t, err, "读取部分文件应该成功")
	assert.NotContains(t, result2, "1 | ", "不应该包含第1行")
	assert.NotContains(t, result2, "2 | ", "不应该包含第2行")
	assert.Contains(t, result2, "3 | Line 3", "应该包含第3行")
	assert.Contains(t, result2, "5 | Line 5", "应该包含第5行")

	// 测试用例3：读取部分文件（指定结束行）
	result3, err := tool.ReadFile("test_file.txt", 0, 2)
	assert.NoError(t, err, "读取部分文件应该成功")
	assert.Contains(t, result3, "1 | Line 1", "应该包含第1行")
	assert.Contains(t, result3, "2 | Line 2", "应该包含第2行")
	assert.NotContains(t, result3, "3 | ", "不应该包含第3行")

	// 测试用例4：读取部分文件（指定开始行和结束行）
	result4, err := tool.ReadFile("test_file.txt", 2, 4)
	assert.NoError(t, err, "读取部分文件应该成功")
	assert.NotContains(t, result4, "1 | ", "不应该包含第1行")
	assert.Contains(t, result4, "2 | Line 2", "应该包含第2行")
	assert.Contains(t, result4, "4 | Line 4", "应该包含第4行")
	assert.NotContains(t, result4, "5 | ", "不应该包含第5行")

	// 测试用例5：读取不存在的文件
	_, err = tool.ReadFile("nonexistent_file.txt", 0, 0)
	assert.Error(t, err, "读取不存在的文件应该返回错误")

	// 测试用例6：开始行大于文件行数
	result6, err := tool.ReadFile("test_file.txt", 10, 0)
	assert.NoError(t, err, "开始行大于文件行数应该不返回错误")
	assert.Empty(t, strings.TrimSpace(result6), "开始行大于文件行数应该返回空内容")

	// 测试用例7：结束行小于开始行
	result7, err := tool.ReadFile("test_file.txt", 3, 2)
	assert.NoError(t, err, "结束行小于开始行应该不返回错误")
	assert.Empty(t, strings.TrimSpace(result7), "结束行小于开始行应该返回空内容")
}
