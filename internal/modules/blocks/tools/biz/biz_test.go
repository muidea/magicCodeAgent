package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParamContent(t *testing.T) {
	// 测试用例1：readFile工具 - 正常XML格式，包含多个参数
	content1 := `<readFile>
	<path>test.txt</path>
	<start_line>1</start_line>
	<end_line>10</end_line>
</readFile>`
	funcName1, params1, err := parseParamContent(content1)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "readFile", funcName1, "函数名应该是readFile")
	assert.Equal(t, "test.txt", params1["path"], "path参数应该正确解析")
	assert.Equal(t, "1", params1["start_line"], "start_line参数应该正确解析")
	assert.Equal(t, "10", params1["end_line"], "end_line参数应该正确解析")

	// 测试用例2：listFiles工具 - 正常XML格式，包含递归参数
	content2 := `<listFiles>
	<path>src</path>
	<recursive>true</recursive>
</listFiles>`
	funcName2, params2, err := parseParamContent(content2)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "listFiles", funcName2, "函数名应该是listFiles")
	assert.Equal(t, "src", params2["path"], "path参数应该正确解析")
	assert.Equal(t, "true", params2["recursive"], "recursive参数应该正确解析")

	// 测试用例3：writeToFile工具 - 正常XML格式，包含多行内容
	content3 := `<writeToFile>
	<path>output.txt</path>
	<content>
This is line 1
This is line 2
This is line 3
	</content>
	<line_count>3</line_count>
</writeToFile>`
	funcName3, params3, err := parseParamContent(content3)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "writeToFile", funcName3, "函数名应该是writeToFile")
	assert.Equal(t, "output.txt", params3["path"], "path参数应该正确解析")
	assert.Equal(t, "This is line 1\nThis is line 2\nThis is line 3\n\t", params3["content"], "content参数应该正确解析")
	assert.Equal(t, "3", params3["line_count"], "line_count参数应该正确解析")

	// 测试用例4：searchFiles工具 - 正常XML格式，包含可选参数
	content4 := `<searchFiles>
	<path>.</path>
	<regex>.*</regex>
	<file_pattern>*.go</file_pattern>
</searchFiles>`
	funcName4, params4, err := parseParamContent(content4)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "searchFiles", funcName4, "函数名应该是searchFiles")
	assert.Equal(t, ".", params4["path"], "path参数应该正确解析")
	assert.Equal(t, ".*", params4["regex"], "regex参数应该正确解析")
	assert.Equal(t, "*.go", params4["file_pattern"], "file_pattern参数应该正确解析")

	// 测试用例5：executeCommand工具 - 正常XML格式，包含可选的cwd参数
	content5 := `<executeCommand>
	<command>ls -la</command>
	<cwd>/home/user/projects</cwd>
</executeCommand>`
	funcName5, params5, err := parseParamContent(content5)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "executeCommand", funcName5, "函数名应该是executeCommand")
	assert.Equal(t, "ls -la", params5["command"], "command参数应该正确解析")
	assert.Equal(t, "/home/user/projects", params5["cwd"], "cwd参数应该正确解析")

	// 测试用例6：searchAndReplace工具 - 正常XML格式，包含复杂JSON参数
	content6 := `<searchAndReplace>
	<path>example.ts</path>
	<operations>[{"search":"foo","replace":"bar","start_line":1,"end_line":10}]</operations>
</searchAndReplace>`
	funcName6, params6, err := parseParamContent(content6)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "searchAndReplace", funcName6, "函数名应该是searchAndReplace")
	assert.Equal(t, "example.ts", params6["path"], "path参数应该正确解析")
	assert.Equal(t, "[{\"search\":\"foo\",\"replace\":\"bar\",\"start_line\":1,\"end_line\":10}]", params6["operations"], "operations参数应该正确解析")

	// 测试用例7：accessMcpResource工具 - 正常XML格式
	content7 := `<access_mcp_resource>
	<server_name>weather-server</server_name>
	<uri>weather://san-francisco/current</uri>
</access_mcp_resource>`
	funcName7, params7, err := parseParamContent(content7)
	assert.NoError(t, err, "解析有效XML应该成功")
	assert.Equal(t, "access_mcp_resource", funcName7, "函数名应该是access_mcp_resource")
	assert.Equal(t, "weather-server", params7["server_name"], "server_name参数应该正确解析")
	assert.Equal(t, "weather://san-francisco/current", params7["uri"], "uri参数应该正确解析")

	// 测试用例8：空XML内容
	content8 := ""
	_, _, err = parseParamContent(content8)
	assert.Error(t, err, "解析空XML内容应该返回错误")

	// 测试用例9：无效XML格式
	content9 := "<invalidXml>"
	_, _, err = parseParamContent(content9)
	assert.Error(t, err, "解析无效XML格式应该返回错误")

	// 测试用例10：XML格式正确但没有root节点
	content10 := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	_, _, err = parseParamContent(content10)
	assert.Error(t, err, "解析没有root节点的XML应该返回错误")

	// 测试用例11：XML格式正确但参数值为空
	content11 := `<writeToFile>
	<path></path>
	<content></content>
	<line_count></line_count>
</writeToFile>`
	funcName11, params11, err := parseParamContent(content11)
	assert.NoError(t, err, "解析参数值为空的XML应该成功")
	assert.Equal(t, "writeToFile", funcName11, "函数名应该是writeToFile")
	assert.Equal(t, "", params11["path"], "空path参数应该正确解析")
	assert.Equal(t, "", params11["content"], "空content参数应该正确解析")
	assert.Equal(t, "", params11["line_count"], "空line_count参数应该正确解析")

	content12 := `<apply_diff>
<path>Makefile</path>
<diff>
<<<<<<< SEARCH
:start_line:8
-------
run: build
	./hello
=======
run: build
	./hello

clean:
	rm -f hello
>>>>>>> REPLACE
</diff>
</apply_diff>`
	funcName12, param12, err := parseParamContent(content12)
	assert.NoError(t, err, "解析失败")
	assert.Equal(t, "apply_diff", funcName12, "函数名应该正确解析")
	assert.Equal(t, "Makefile", param12["path"], "path参数应该正确解析")

	content13 := `<write_to_file>
<path>product_list.md</path>
<content>
# 产品列表页面

## 功能描述
- 展示所有产品的缩略信息
- 支持分页加载
- 支持下拉刷新

## 接口需求
- GET /api/products?page=1&limit=10

## 交互流程
1. 进入页面自动加载第一页数据
2. 滚动到底部自动加载下一页
3. 下拉刷新重置列表

## UI组件
- 产品卡片组件
- 无限滚动组件
- 下拉刷新组件
- 搜索框组件
</content>
<line_count>16</line_count>
</write_to_file>`
	funcName13, param13, err := parseParamContent(content13)
	assert.NoError(t, err, "解析失败")
	assert.Equal(t, "write_to_file", funcName13, "函数名应该正确解析")
	assert.Equal(t, "product_list.md", param13["path"], "path参数应该正确解析")

	content14 := `<ask_followup_question>
<question>abc</question>
<follow_up>
<suggest>需要补充会员签到和优惠券功能</suggest>
<suggest>需要补充会员分组/标签功能</suggest>
<suggest>需要补充会员活动报名功能</suggest>
<suggest>当前分解已完整，无需补充</suggest>
</follow_up>
</ask_followup_question>`
	funcName14, param14, err := parseParamContent(content14)
	assert.NoError(t, err, "parseParamContent error")
	assert.Equal(t, "ask_followup_question", funcName14, "funcName14 error")
	assert.Equal(t, "abc", param14["question"], "param14[question] error")
	assert.Equal(t, 4, len(param14["follow_up"].([]string)), "param14[follow_up.suggest[0]] error")

	content15 := `<execute_command>
<command>mkdir -p src/pages/auth src/services src/models</command>
</execute_command>`
	funcName15, param15, err := parseParamContent(content15)
	assert.NoError(t, err, "parseParamContent error")
	assert.Equal(t, "execute_command", funcName15, "funcName15 error")
	assert.Equal(t, 1, len(param15), "param15 error")
}
