package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchAndReplace(t *testing.T) {
	// 临时修改WorkspaceDir为测试目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFilePath := filepath.Join(tempDir, "test_file.txt")
	testContent := "Line 1: Hello World\nLine 2: Test Content\nLine 3: Replace Me\nLine 4: Test Content\nLine 5: Final Line"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")

	// 测试用例1：简单字符串替换
	op1 := []struct {
		Search     string `json:"search"`
		Replace    string `json:"replace"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		UseRegex   bool   `json:"use_regex"`
		IgnoreCase bool   `json:"ignore_case"`
		RegexFlags string `json:"regex_flags"`
	}{
		{
			Search:  "Hello World",
			Replace: "Hello Universe",
		},
	}

	err = SearchAndReplace(tempDir, "test_file.txt", op1)
	assert.NoError(t, err, "简单字符串替换应该成功")

	// 验证替换结果
	content, err := os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")
	assert.Contains(t, string(content), "Line 1: Hello Universe", "应该替换Hello World为Hello Universe")
	assert.Contains(t, string(content), "Line 3: Replace Me", "不应该修改其他行")

	// 重置测试文件
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "重置测试文件应该成功")

	// 测试用例2：指定行范围的替换
	op2 := []struct {
		Search     string `json:"search"`
		Replace    string `json:"replace"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		UseRegex   bool   `json:"use_regex"`
		IgnoreCase bool   `json:"ignore_case"`
		RegexFlags string `json:"regex_flags"`
	}{
		{
			Search:    "Test Content",
			Replace:   "Modified Content",
			StartLine: 2,
			EndLine:   3,
		},
	}

	err = SearchAndReplace(tempDir, "test_file.txt", op2)
	assert.NoError(t, err, "指定行范围的替换应该成功")

	// 验证替换结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")
	assert.Contains(t, string(content), "Line 2: Modified Content", "应该在指定行范围内替换")
	assert.Contains(t, string(content), "Line 4: Test Content", "不应该修改范围外的行")

	// 重置测试文件
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "重置测试文件应该成功")

	// 测试用例3：使用正则表达式替换
	op3 := []struct {
		Search     string `json:"search"`
		Replace    string `json:"replace"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		UseRegex   bool   `json:"use_regex"`
		IgnoreCase bool   `json:"ignore_case"`
		RegexFlags string `json:"regex_flags"`
	}{
		{
			Search:    "Line \\d+: (.*)",
			Replace:   "Modified $1",
			UseRegex:  true,
			StartLine: 3,
			EndLine:   4,
		},
	}

	err = SearchAndReplace(tempDir, "test_file.txt", op3)
	assert.NoError(t, err, "使用正则表达式替换应该成功")

	// 验证替换结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")
	assert.Contains(t, string(content), "Modified Replace Me", "应该使用正则表达式替换")
	assert.Contains(t, string(content), "Line 1: Hello World", "不应该修改范围外的行")

	// 重置测试文件
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(t, err, "重置测试文件应该成功")

	// 测试用例4：多个替换操作
	op4 := []struct {
		Search     string `json:"search"`
		Replace    string `json:"replace"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		UseRegex   bool   `json:"use_regex"`
		IgnoreCase bool   `json:"ignore_case"`
		RegexFlags string `json:"regex_flags"`
	}{
		{
			Search:  "Hello World",
			Replace: "Hello Universe",
		},
		{
			Search:  "Replace Me",
			Replace: "Replaced Successfully",
		},
	}

	err = SearchAndReplace(tempDir, "test_file.txt", op4)
	assert.NoError(t, err, "多个替换操作应该成功")

	// 验证替换结果
	content, err = os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")
	assert.Contains(t, string(content), "Line 1: Hello Universe", "应该替换Hello World为Hello Universe")
	assert.Contains(t, string(content), "Line 3: Replaced Successfully", "应该替换Replace Me为Replaced Successfully")

	// 测试用例5：替换不存在的文件
	op5 := []struct {
		Search     string `json:"search"`
		Replace    string `json:"replace"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		UseRegex   bool   `json:"use_regex"`
		IgnoreCase bool   `json:"ignore_case"`
		RegexFlags string `json:"regex_flags"`
	}{
		{
			Search:  "test",
			Replace: "replaced",
		},
	}

	err = SearchAndReplace(tempDir, "nonexistent_file.txt", op5)
	assert.Error(t, err, "替换不存在的文件应该返回错误")
}
