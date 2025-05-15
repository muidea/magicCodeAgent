package biz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type ExecuteCommandTool struct{}

func (s *ExecuteCommandTool) Name() string { return "execute_command" }
func (s *ExecuteCommandTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ExecuteCommandTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Result: models.MessageItem{
				Content: []string{
					fmt.Sprintf("execute tool: %s, param:%v", s.Name(), param),
				},
			},
		}
		if err != nil {
			ret.Final = true
			ret.Result.Content = append(ret.Result.Content, err.Error())
		} else {
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	command, ok := param["command"].(string)
	if !ok || command == "" {
		err = fmt.Errorf("command is empty")
		return
	}
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	var args []string
	if argsParam, ok := param["args"].([]interface{}); ok {
		for _, arg := range argsParam {
			if strArg, ok := arg.(string); ok {
				args = append(args, strArg)
			}
		}
	}

	cwd, _ := param["cwd"].(string)
	result, err = ExecuteCommand(fullPath, command, args, cwd)
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}

/*
## execute_command
Description: Request to execute a CLI command on the system. Use this when you need to perform system operations or run specific commands to accomplish any step in the user's task. You must tailor your command to the user's system and provide a clear explanation of what the command does. For command chaining, use the appropriate chaining syntax for the user's shell. Prefer to execute complex CLI commands over creating executable scripts, as they are more flexible and easier to run. Prefer relative commands and paths that avoid location sensitivity for terminal consistency, e.g: `touch ./testdata/example.file`, `dir ./examples/model1/data/yaml`, or `go test ./cmd/front --config ./cmd/front/config.yml`. If directed by the user, you may open a terminal in a different directory by using the `cwd` parameter.
描述：请求在系统上执行CLI命令。当需要执行系统操作或运行特定命令来完成用户任务的任何步骤时使用此功能。必须根据用户的系统定制命令，并清晰解释命令的作用。对于命令链式操作，使用适合用户shell的链式语法。优先执行复杂的CLI命令而非创建可执行脚本，因其更灵活且易于运行。优先使用相对命令和路径以确保终端一致性，例如：`touch ./testdata/example.file`、`dir ./examples/model1/data/yaml` 或 `go test ./cmd/front --config ./cmd/front/config.yml`。如果用户指定，可通过`cwd`参数在不同目录打开终端。
Parameters:
- command: (required) The CLI command to execute. This should be valid for the current operating system. Ensure the command is properly formatted and does not contain any harmful instructions.
  命令：(必填) 要执行的CLI命令，必须适用于当前操作系统。确保命令格式正确且不包含有害指令。
- cwd: (optional) The working directory to execute the command in (default: /home/rangh/webspace)
  工作目录：(可选) 执行命令的工作目录（默认：/home/rangh/webspace）
Usage:
<execute_command>
<command>Your command here</command>
<cwd>Working directory path (optional)</cwd>
</execute_command>

Example: Requesting to execute npm run dev
示例：请求执行npm run dev
<execute_command>
<command>npm run dev</command>
</execute_command>

Example: Requesting to execute ls in a specific directory if directed
示例：按要求在指定目录执行ls
<execute_command>
<command>ls -la</command>
<cwd>/home/user/projects</cwd>
</execute_command>

*/

// ExecuteCommand 实现execute_command工具功能
func ExecuteCommand(workspaceDir, command string, args []string, cwd string) (string, error) {
	items := strings.Split(command, " ")
	items = append(items, args...)
	if len(items) > 0 {
		command = items[0]
		args = items[1:]
	}
	cmd := exec.Command(command, args...)

	if cwd != "" {
		cmd.Dir = filepath.Join(workspaceDir, cwd)
	} else {
		cmd.Dir = workspaceDir
	}

	if !fp.Exist(cmd.Dir) {
		if err := os.MkdirAll(cmd.Dir, 0755); err != nil {
			log.Errorf("failed to create directory: %v, err：%s", cmd.Dir, err)
			return "", err
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("failed to execute command: %v, err：%s", cmd, err)
		return string(output), err
	}

	return string(output), nil
}
