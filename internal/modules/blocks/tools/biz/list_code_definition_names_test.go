package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestListCodeDefinitionNames(t *testing.T) {
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
	testFilePath := filepath.Join(tempDir, "test_file.go")
	testContent := `package test

func TestFunction() string {
	return "test"
}

type TestStruct struct {
	Field string
}

var TestVar = "test"
`
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")

	// 测试用例1：分析单个文件
	definitions, err := tool.ListCodeDefinitionNames("test_file.go")
	assert.NoError(t, err, "分析单个文件应该成功")
	assert.Contains(t, definitions, "TestFunction", "应该包含函数定义")
	assert.Contains(t, definitions, "TestStruct", "应该包含结构体定义")
	assert.Contains(t, definitions, "TestVar", "应该包含变量定义")

	// 创建测试目录和多个文件
	testDirPath := filepath.Join(tempDir, "testdir")
	err = os.Mkdir(testDirPath, 0755)
	assert.NoError(t, err, "创建测试目录应该成功")

	// 创建第一个测试文件
	testFile1Path := filepath.Join(testDirPath, "file1.go")
	testContent1 := `package testdir

func File1Function() {}

type File1Struct struct {}
`
	err = os.WriteFile(testFile1Path, []byte(testContent1), 0644)
	assert.NoError(t, err, "写入第一个测试文件应该成功")

	// 创建第二个测试文件
	testFile2Path := filepath.Join(testDirPath, "file2.go")
	testContent2 := `package testdir

func File2Function() {}

var File2Var = 42
`
	err = os.WriteFile(testFile2Path, []byte(testContent2), 0644)
	assert.NoError(t, err, "写入第二个测试文件应该成功")

	// 测试用例2：分析目录
	definitions2, err := tool.ListCodeDefinitionNames("testdir")
	assert.NoError(t, err, "分析目录应该成功")
	assert.Contains(t, definitions2, "File1Function", "应该包含第一个文件中的函数定义")
	assert.Contains(t, definitions2, "File1Struct", "应该包含第一个文件中的结构体定义")
	assert.Contains(t, definitions2, "File2Function", "应该包含第二个文件中的函数定义")
	assert.Contains(t, definitions2, "File2Var", "应该包含第二个文件中的变量定义")

	// 测试用例3：分析不存在的路径
	_, err = tool.ListCodeDefinitionNames("nonexistent_path")
	assert.Error(t, err, "分析不存在的路径应该返回错误")

	// 测试用例4：分析非Go文件
	nonGoFilePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(nonGoFilePath, []byte("This is not a Go file"), 0644)
	assert.NoError(t, err, "写入非Go文件应该成功")

	_, err = tool.ListCodeDefinitionNames("test.txt")
	assert.Error(t, err, "分析非Go文件应该返回错误")
}
