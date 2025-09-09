package models

import "testing"

func TestParseTodo(t *testing.T) {
	content := `[x] 移除搜索输入框和相关状态
[x] 修改renderCardView函数，在List中添加新增卡片
[-] 移除未使用的导入
[ ] 修复条件判断逻辑
[ ] 调整样式和逻辑以支持新增卡片`

	todoList := ParseTodoList(content)
	for _, todo := range todoList {
		t.Logf("Task: %s, Status: %s", todo.Content, todo.Status)
	}
	if len(todoList) != 5 {
		t.Errorf("Expected 5 tasks, got %d", len(todoList))
	}
}
