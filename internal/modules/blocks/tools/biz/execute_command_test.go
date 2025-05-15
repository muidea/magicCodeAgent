package biz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommand(t *testing.T) {
	// 临时修改WorkspaceDir为测试目录
	tempDir := t.TempDir()
	// 创建测试文件
	testFilePath := filepath.Join(tempDir, "test_file.txt")
	testContent := "Hello, World!"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")

	// 测试用例1：执行ls命令
	output1, err := ExecuteCommand(tempDir, "ls", []string{"-la"}, "")
	assert.NoError(t, err, "执行ls命令应该成功")
	assert.Contains(t, output1, "test_file.txt", "输出应该包含测试文件名")

	// 测试用例2：执行cat命令读取文件内容
	output2, err := ExecuteCommand(tempDir, "cat", []string{"test_file.txt"}, "")
	assert.NoError(t, err, "执行cat命令应该成功")
	assert.Equal(t, testContent, strings.TrimSpace(output2), "输出应该是文件内容")

	// 测试用例3：在指定目录执行命令
	// 创建子目录
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err, "创建子目录应该成功")

	// 在子目录中创建文件
	subFilePath := filepath.Join(subDir, "sub_file.txt")
	subContent := "Subdirectory file"
	err = os.WriteFile(subFilePath, []byte(subContent), 0644)
	assert.NoError(t, err, "写入子目录测试文件应该成功")

	// 在子目录中执行ls命令
	output3, err := ExecuteCommand(tempDir, "ls", []string{"-la"}, "subdir")
	assert.NoError(t, err, "在子目录执行ls命令应该成功")
	assert.Contains(t, output3, "sub_file.txt", "输出应该包含子目录测试文件名")

	// 测试用例4：执行不存在的命令
	_, err = ExecuteCommand(tempDir, "nonexistent_command", []string{}, "")
	assert.Error(t, err, "执行不存在的命令应该返回错误")

	_, err = ExecuteCommand(tempDir, "ls -la", []string{}, "")
	assert.NoError(t, err, "在根目录执行ls命令应该成功")

	_, err = ExecuteCommand(tempDir, "mkdir -p src/pages/auth src/services src/models", []string{}, "")
	assert.NoError(t, err, "在根目录创建子目录应该成功")
}
