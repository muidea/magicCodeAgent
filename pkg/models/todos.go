package models

import (
	"fmt"
	"strings"
)

type TodoItem struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

type TodoList []*TodoItem

/*
格式化成如下格式
| # | Content | Status |
|---|---------|--------|
| 1 | 移除搜索输入框和相关状态 | Completed |
| 2 | 修改renderCardView函数，在List中添加新增卡片 | Completed |
| 3 | 移除未使用的导入 | Pending |
| 4 | 修复条件判断逻辑 | Pending |
| 5 | 调整样式和逻辑以支持新增卡片 | Pending |
*/
func (s *TodoList) String() string {
	if len(*s) == 0 {
		return "| # | Content | Status |\n|---|---------|--------|"
	}

	var result string
	result += "| # | Content | Status |\n"
	result += "|---|---------|--------|\n"

	for i, item := range *s {
		result += fmt.Sprintf("| %d | %s | %s |\n", i+1, item.Content, item.Status)
	}

	return result
}

/*
解析如下格式的内容，以生成一个 TodoList
规则如下:
- [ ] 任务描述（待处理）
- [x] 任务描述（已完成）
- [-] 任务描述（进行中）
-----
[x] 移除搜索输入框和相关状态
[x] 修改renderCardView函数，在List中添加新增卡片
[-] 移除未使用的导入
[ ] 修复条件判断逻辑
[ ] 调整样式和逻辑以支持新增卡片
*/
func ParseTodoList(content string) (todoList TodoList) {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行和分隔线
		if line == "" {
			continue
		}

		// 解析任务行格式: [状态] 任务描述
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			// 提取状态标记
			endBracket := strings.Index(line, "]")
			if endBracket > 0 {
				statusMark := strings.TrimSpace(line[1:endBracket])
				taskContent := strings.TrimSpace(line[endBracket+1:])

				// 映射状态标记到状态字符串
				var status string
				switch statusMark {
				case "x":
					status = "Completed"
				case "-":
					status = "In Progress"
				case "":
					status = "Pending"
				default:
					// 未知状态标记，跳过或使用默认状态
					continue
				}

				// 添加到待办列表
				todoList = append(todoList, &TodoItem{
					Content: taskContent,
					Status:  status,
				})
			}
		}
	}

	return todoList
}
