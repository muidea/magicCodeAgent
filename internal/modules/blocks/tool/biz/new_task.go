package biz

import "fmt"

/*
## new_task
Description: Create a new task with a specified starting mode and initial message. This tool instructs the system to create a new Cline instance in the given mode with the provided message.

Parameters:
- mode: (required) The slug of the mode to start the new task in (e.g., \"code\", \"ask\", \"architect\").
- message: (required) The initial user message or instructions for this new task.

Usage:
<new_task>
<mode>your-mode-slug-here</mode>
<message>Your initial instructions here</message>
</new_task>

Example:
<new_task>
<mode>code</mode>
<message>Implement a new feature for the application.</message>
</new_task>
*/

// NewTask 实现new_task工具功能
func (s *Tool) NewTask(mode, message string) (string, error) {
	// 返回新任务创建信息
	return fmt.Sprintf("Created new task in mode '%s' with message: %s", mode, message), nil
}
