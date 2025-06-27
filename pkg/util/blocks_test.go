package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testVal = `<thinking>
1. 分析任务：需要为已有的Go程序添加Makefile，用于构建和运行
2. 目标：创建一个标准的Makefile，包含build和run目标
3. 工具选择：使用write_to_file工具创建新的Makefile
4. 参数确认：
   - path: 命名为Makefile（标准名称）
   - content: 包含build和run规则的Makefile内容
   - line_count: 内容共8行（包括注释和空行）
5. 所有必要参数都已明确，可以执行
</thinking>

<write_to_file>
<path>Makefile</path>
<content>
# Makefile for Go demo

.PHONY: build run

build:
	go build -o hello main.go

run:
	go run main.go
</content>
<line_count>8</line_count>
</write_to_file>

<attempt_completion>
<result>
已创建完整的Go演示项目，包含：
1. main.go - 输出"hello world"的Go程序
2. Makefile - 提供build和run命令的构建文件

您可以使用以下命令运行程序：
- make build 编译程序
- make run 直接运行程序
- ./hello 运行编译后的程序
</result>
<command>make run</command>
</attempt_completion>
`

// 空输入测试用例
var emptyInput = ``

// 嵌套标签测试用例
var nestedTags = `
<outer>
这是外部标签的内容
<inner>
这是内部嵌套的标签内容
</inner>
外部标签的结束内容
</outer>
`

// 不匹配标签测试用例
var mismatchedTags = `
<tag1>
这是一个标签的内容
</tag2>
`

// 特殊字符测试用例
var specialChars = `
<special_chars>
特殊字符：!@#$%^&*()_+{}|:"<>?[]\
中文内容：你好，世界！
换行符和制表符：\n\t
</special_chars>
`

// 多个连续标签测试用例
var multipleTags = `
<tag1>
第一个标签的内容
</tag1>
<tag2>
第二个标签的内容
</tag2>
<tag3>
第三个标签的内容
</tag3>
`

func TestBlocks(t *testing.T) {
	// 测试正常输入
	t.Run("正常输入", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(testVal)
		assert.Nil(t, blocksErr)
		assert.Equal(t, 3, len(blocksList))
		assert.Contains(t, blocksList[0], "<thinking>")
		assert.Contains(t, blocksList[1], "<write_to_file>")
		assert.Contains(t, blocksList[2], "<attempt_completion>")
	})

	// 测试空输入
	t.Run("空输入", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(emptyInput)
		assert.Nil(t, blocksErr)
		assert.Equal(t, 0, len(blocksList))
	})

	// 测试嵌套标签
	t.Run("嵌套标签", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(nestedTags)
		assert.Nil(t, blocksErr)
		assert.Equal(t, 1, len(blocksList))
		assert.Contains(t, blocksList[0], "<outer>")
		assert.Contains(t, blocksList[0], "<inner>")
		assert.Contains(t, blocksList[0], "</inner>")
		assert.Contains(t, blocksList[0], "</outer>")
	})

	// 测试不匹配标签
	t.Run("不匹配标签", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(mismatchedTags)
		assert.Nil(t, blocksErr)
		// 由于结束标签不匹配，应该没有完整的块
		assert.Equal(t, 0, len(blocksList))
	})

	// 测试特殊字符
	t.Run("特殊字符", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(specialChars)
		assert.Nil(t, blocksErr)
		assert.Equal(t, 1, len(blocksList))
		assert.Contains(t, blocksList[0], "<special_chars>")
		assert.Contains(t, blocksList[0], "特殊字符")
		assert.Contains(t, blocksList[0], "中文内容")
		assert.Contains(t, blocksList[0], "</special_chars>")
	})

	// 测试多个连续标签
	t.Run("多个连续标签", func(t *testing.T) {
		blocksList, blocksErr := SpliteBlocks(multipleTags)
		assert.Nil(t, blocksErr)
		assert.Equal(t, 3, len(blocksList))
		assert.Contains(t, blocksList[0], "<tag1>")
		assert.Contains(t, blocksList[1], "<tag2>")
		assert.Contains(t, blocksList[2], "<tag3>")
	})
}
