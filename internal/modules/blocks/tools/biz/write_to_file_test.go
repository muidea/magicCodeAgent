package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	// 临时修改WorkspaceDir为测试目录
	tempDir := t.TempDir()

	// 测试用例1：写入新文件
	filePath := "test_file.txt"
	content := "Hello, World!\nThis is a test file.\nLine 3"
	lineCount := 3

	err := WriteToFile(tempDir, filePath, content, lineCount)
	assert.NoError(t, err, "写入新文件应该成功")

	// 验证文件内容
	fullPath := filepath.Join(tempDir, filePath)
	readContent, err := os.ReadFile(fullPath)
	assert.NoError(t, err, "读取文件应该成功")
	assert.Equal(t, content, string(readContent), "文件内容应该与写入的内容一致")

	// 测试用例2：覆盖现有文件
	newContent := "This is new content.\nReplacing the old content."
	newLineCount := 2

	err = WriteToFile(tempDir, filePath, newContent, newLineCount)
	assert.NoError(t, err, "覆盖现有文件应该成功")

	// 验证文件内容已更新
	readContent, err = os.ReadFile(fullPath)
	assert.NoError(t, err, "读取文件应该成功")
	assert.Equal(t, newContent, string(readContent), "文件内容应该已更新")

	// 测试用例3：写入嵌套目录中的文件
	nestedFilePath := "nested/dir/test_file.txt"
	nestedContent := "Content in nested directory"
	nestedLineCount := 1

	err = WriteToFile(tempDir, nestedFilePath, nestedContent, nestedLineCount)
	assert.NoError(t, err, "写入嵌套目录中的文件应该成功")

	// 验证嵌套目录和文件已创建
	fullNestedPath := filepath.Join(tempDir, "nested/dir/test_file.txt")
	readContent, err = os.ReadFile(fullNestedPath)
	assert.NoError(t, err, "读取嵌套目录中的文件应该成功")
	assert.Equal(t, nestedContent, string(readContent), "嵌套目录中的文件内容应该正确")

	// 测试用例4：写入空内容
	emptyFilePath := "empty_file.txt"
	emptyContent := ""
	emptyLineCount := 0

	err = WriteToFile(tempDir, emptyFilePath, emptyContent, emptyLineCount)
	assert.NoError(t, err, "写入空内容应该成功")

	// 验证空文件已创建
	fullEmptyPath := filepath.Join(tempDir, emptyFilePath)
	readContent, err = os.ReadFile(fullEmptyPath)
	assert.NoError(t, err, "读取空文件应该成功")
	assert.Equal(t, emptyContent, string(readContent), "空文件内容应该为空")

	// 测试用例5：写入包含特殊字符的内容
	specialFilePath := "special_file.txt"
	specialContent := "Line with special chars: !@#$%^&*()\nUnicode: 你好，世界\nTab:\tIndented"
	specialLineCount := 3

	err = WriteToFile(tempDir, specialFilePath, specialContent, specialLineCount)
	assert.NoError(t, err, "写入包含特殊字符的内容应该成功")

	// 验证特殊字符文件内容
	fullSpecialPath := filepath.Join(tempDir, specialFilePath)
	readContent, err = os.ReadFile(fullSpecialPath)
	assert.NoError(t, err, "读取包含特殊字符的文件应该成功")
	assert.Equal(t, specialContent, string(readContent), "包含特殊字符的文件内容应该正确")
}
