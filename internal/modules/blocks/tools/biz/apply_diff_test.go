package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitDiff(t *testing.T) {
	tests := []struct {
		name        string
		diff        string
		expected    []*DiffBlock
		expectedErr bool
	}{
		{
			name: "缺少end_line的diff块",
			diff: `<<<<<<< SEARCH
:start_line:1
-------
line1
line2
line3
=======
new line1
new line2
new line3
>>>>>>> REPLACE`,
			expected: []*DiffBlock{
				{
					StartLine:      1,
					EndLine:        3, // 应该根据内容行数自动计算为3
					OriginContent:  "line1\nline2\nline3",
					ReplaceContent: "new line1\nnew line2\nnew line3",
				},
			},
		},
		{
			name: "单个diff块",
			diff: `<<<<<<< SEARCH
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
    """Calculate total with 10% markup"""
    return sum(item * 1.1 for item in items)
>>>>>>> REPLACE`,
			expected: []*DiffBlock{
				{
					StartLine:      1,
					EndLine:        5,
					OriginContent:  "def calculate_total(items):\n    total = 0\n    for item in items:\n        total += item\n    return total",
					ReplaceContent: "def calculate_total(items):\n    \"\"\"Calculate total with 10% markup\"\"\"\n    return sum(item * 1.1 for item in items)",
				},
			},
		},
		{
			name: "多个diff块",
			diff: `<<<<<<< SEARCH
:start_line:10
:end_line:12
-------
func foo() {
    // do something
}
=======
func foo() {
    // do something better
}
>>>>>>> REPLACE

<<<<<<< SEARCH
:start_line:5
:end_line:7
-------
func init() {
    // initialize
}
=======
func init() {
    // better initialization
}
>>>>>>> REPLACE`,
			expected: []*DiffBlock{
				{
					StartLine:      5,
					EndLine:        7,
					OriginContent:  "func init() {\n    // initialize\n}",
					ReplaceContent: "func init() {\n    // better initialization\n}",
				},
				{
					StartLine:      10,
					EndLine:        12,
					OriginContent:  "func foo() {\n    // do something\n}",
					ReplaceContent: "func foo() {\n    // do something better\n}",
				},
			},
		},
		{
			name: "格式不正确的diff",
			diff: `<<<<<<< SEARCH
:start_line:1
:end_line:3
-------
func bar() {
=======
func newBar() {
>>>>>>> REPLACE`,
			expected: []*DiffBlock{
				{
					StartLine:      1,
					EndLine:        3,
					OriginContent:  "func bar() {",
					ReplaceContent: "func newBar() {",
				},
			},
			expectedErr: false,
		},
		{
			name: "起始行大于结束行",
			diff: `<<<<<<< SEARCH
:start_line:5
:end_line:3
-------
func bar() {
=======
func newBar() {
>>>>>>> REPLACE`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "缺少分隔符",
			diff: `<<<<<<< SEARCH
:start_line:1
:end_line:3
func bar() {
=======
func newBar() {
>>>>>>> REPLACE`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "空diff",
			diff:        "",
			expected:    nil,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := splitDiffBlock(tt.diff)

			// 检查错误情况
			if (err != nil) != tt.expectedErr {
				t.Errorf("Expected error: %v, got: %v, error: %v", tt.expectedErr, err != nil, err)
				return
			}

			if err != nil {
				// 如果期望有错误，并且确实有错误，测试通过
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d blocks, got %d", len(tt.expected), len(result))
				return
			}

			for i, block := range result {
				expected := tt.expected[i]
				if block.StartLine != expected.StartLine ||
					block.EndLine != expected.EndLine ||
					block.OriginContent != expected.OriginContent ||
					block.ReplaceContent != expected.ReplaceContent {
					t.Errorf("Block %d mismatch:\nExpected: %+v\nGot: %+v",
						i, expected, block)
				}
			}
		})
	}
}

func TestApplyDiff(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test_file.txt")

	// 写入初始内容
	initialContent := "def calculate_total(items):\n    total = 0\n    for item in items:\n        total += item\n    return total"
	err := os.WriteFile(testFilePath, []byte(initialContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")

	// 测试用例1：单个替换块
	diff1 := `<<<<<<< SEARCH
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
    """Calculate total with 10% markup"""
    return sum(item * 1.1 for item in items)
>>>>>>> REPLACE
`

	// 应用差异
	err = ApplyDiff(tempDir, "test_file.txt", diff1)
	assert.NoError(t, err, "应用差异应该成功")

	// 读取修改后的文件内容
	modifiedContent, err := os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")

	// 验证修改结果
	expectedContent := "def calculate_total(items):\n    \"\"\"Calculate total with 10% markup\"\"\"\n    return sum(item * 1.1 for item in items)"
	assert.Equal(t, expectedContent, string(modifiedContent), "文件内容应该被正确修改")

	// 测试用例2：多个替换块
	// 重置文件内容
	initialContent2 := "def calculate_total(items):\n    sum = 0\n    for item in items:\n        total += item\n    return total"
	err = os.WriteFile(testFilePath, []byte(initialContent2), 0644)
	assert.NoError(t, err, "重置测试文件应该成功")

	diff2 := `<<<<<<< SEARCH
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
`

	// 应用差异
	err = ApplyDiff(tempDir, "test_file.txt", diff2)
	assert.NoError(t, err, "应用多个替换块应该成功")

	// 读取修改后的文件内容
	modifiedContent2, err := os.ReadFile(testFilePath)
	assert.NoError(t, err, "读取修改后的文件应该成功")

	// 验证修改结果
	expectedContent2 := "def calculate_sum(items):\n    sum = 0\n    for item in items:\n        sum += item\n    return sum"
	assert.Equal(t, expectedContent2, string(modifiedContent2), "文件内容应该被正确修改")

	// 测试用例3：文件不存在
	err = ApplyDiff(tempDir, "non_existent_file.txt", diff1)
	assert.Error(t, err, "对不存在的文件应用差异应该返回错误")

	// 测试用例4：无效的差异格式
	invalidDiff := "This is not a valid diff format"
	err = ApplyDiff(tempDir, "test_file.txt", invalidDiff)
	assert.Error(t, err, "应用无效的差异格式应该返回错误")

	diff3 := `<<<<<<< SEARCH
  routes: [
    {
      name: 'index',
      path: '/',
      component: './index',
    },
  ],
=======
  routes: [
    {
      name: 'index',
      path: '/',
      component: './index',
    },
    {
      name: 'member',
      path: '/member',
      component: './member',
      routes: [
        {
          path: '/member/profile',
          component: './member/profile',
        },
        {
          path: '/member/address',
          component: './member/address',
        },
      ],
    },
  ],
>>>>>>> REPLACE`

	blocks, blockErr := splitDiffBlock(diff3)
	assert.Nil(t, blockErr)
	assert.Equal(t, 1, len(blocks))

	rawContent := `import { defineConfig } from '@umijs/max';

export default defineConfig({
  links: [
    { rel: 'icon', href: '/favicon.png', type: 'image/png' },
  ],
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    title: 'demo',
  },
  base: '/',
  publicPath: '/',
  routes: [
    {
      name: 'index',
      path: '/',
      component: './index',
    },
  ],
  exportStatic: {},
  esbuildMinifyIIFE: true,
  npmClient: 'yarn',
  proxy: [
    {
      context: ['/api/v1/'],
      target: 'http://127.0.0.1:8090',
      changeOrigin: true,
    }
  ],
  alias: {
    '@/utils': require('path').resolve(__dirname, 'src/utils'),
    '@/types': require('path').resolve(__dirname, 'src/types'),
  },
});`
	testFilePath = filepath.Join(tempDir, "test_file2.txt")
	err = os.WriteFile(testFilePath, []byte(rawContent), 0644)
	assert.NoError(t, err, "写入测试文件应该成功")
	err = ApplyDiff(tempDir, "test_file2.txt", diff3)
	assert.NoError(t, err, "应用 diff3 应该成功")
}
