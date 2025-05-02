package biz

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertContent(t *testing.T) {
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

	// 测试用例1：在文件开头插入内容
	operationsJSON1 := `[
		{
			"start_line": 1,
			"content": "New Line 0"
		}
	]`

	var operations1 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON1), &operations1)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations1)
	assert.NoError(t, err, "在文件开头插入内容应该成功")

	// 验证插入结果
	content, err := os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取文件应该成功")
	assert.Contains(t, string(content), "New Line 0\nLine 1", "应该在文件开头插入内容")

	// 测试用例2：在文件中间插入多行内容
	operationsJSON2 := `[
		{
			"start_line": 3,
			"content": "New Line A\nNew Line B"
		}
	]`
	var operations2 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON2), &operations2)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations2)
	assert.NoError(t, err, "在文件中间插入多行内容应该成功")

	// 验证插入结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取文件应该成功")
	assert.Contains(t, string(content), "Line 1\nNew Line A\nNew Line B\nLine 2", "应该在文件中间插入多行内容")

	// 测试用例3：在文件末尾插入内容
	operationsJSON3 := `[
		{
			"start_line": 8,
			"content": "New Line End"
		}
	]`
	var operations3 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON3), &operations3)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations3)
	assert.NoError(t, err, "在文件末尾插入内容应该成功")

	// 验证插入结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取文件应该成功")
	assert.True(t, strings.HasSuffix(string(content), "New Line End\nLine 5"), "应该在文件末尾插入内容")

	// 测试用例4：多个位置同时插入内容
	operationsJSON4 := `[
		{
			"start_line": 2,
			"content": "Insert at 2"
		},
		{
			"start_line": 5,
			"content": "Insert at 5"
		}
	]`
	var operations4 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON4), &operations4)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations4)
	assert.NoError(t, err, "多个位置同时插入内容应该成功")

	// 验证插入结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取文件应该成功")
	contentStr := string(content)
	assert.Contains(t, contentStr, "New Line 0\nInsert at 2\nLine 1", "应该在位置2插入内容")
	assert.Contains(t, contentStr, "New Line B\nInsert at 5\nLine 2", "应该在位置5插入内容")

	// 测试用例5：插入位置超出文件范围
	operationsJSON5 := `[
		{
			"start_line": 100,
			"content": "Invalid position"
		}
	]`
	var operations5 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON5), &operations5)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations5)
	assert.Error(t, err, "插入位置超出文件范围应该返回错误")

	// 测试用例6：插入位置为负数
	operationsJSON6 := `[
		{
			"start_line": -1,
			"content": "Invalid position"
		}
	]`
	var operations6 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON6), &operations6)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("test_file.txt", operations6)
	assert.Error(t, err, "插入位置为负数应该返回错误")

	// 测试用例7：文件不存在
	operationsJSON7 := `[
		{
			"start_line": 1,
			"content": "Content"
		}
	]`
	var operations7 []*InsertBlock
	err = json.Unmarshal([]byte(operationsJSON7), &operations7)
	assert.NoError(t, err, "解析JSON应该成功")
	err = InsertContent("nonexistent_file.txt", operations7)
	assert.Error(t, err, "文件不存在应该返回错误")
}
