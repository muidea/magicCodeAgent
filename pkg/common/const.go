package common

import "github.com/sashabaranov/go-openai"

const (
	v1Version      = "/api/v1"
	defaultVersion = v1Version
)

const ApiVersion = defaultVersion

const (
	UrlMsg       = "url"
	SendMsg      = "send"
	RecvMsg      = "recv"
	ErrorMsg     = "error"
	TerminateMsg = "terminate"
)

const (
	System    = openai.ChatMessageRoleSystem
	User      = openai.ChatMessageRoleUser
	Assistant = openai.ChatMessageRoleAssistant
	Function  = openai.ChatMessageRoleFunction
	Tool      = openai.ChatMessageRoleTool
	Developer = openai.ChatMessageRoleDeveloper
	// 以下的自定义角色只是为了 前端展示使用，实际不会用到
	Notifier = "notifier"
)

const (
	UnknownStage = ""
	DesignStage  = "design"
	CodeStage    = "code"
	TestStage    = "test"
	PublishStage = "publish"
)

func ValidStage(stage string) bool {
	switch stage {
	case DesignStage, CodeStage, TestStage, PublishStage:
		return true
	default:
		return false
	}
}

const (
	WorkspaceName  = "workspace"
	WorkspaceDir   = "workspaceDir"
	WorkspaceStage = "workspaceStage"
	WorkspaceType  = "workspaceType"
)

const (
	Service = "service"
	Web     = "web"
)

const (
	Project = "project"
	Module  = "module"
	Block   = "block"
	Point   = "point"
)

const (
	LevelGreater = 1
	LevelEqual   = 0
	LevelLess    = -1
)

// 比较l1, l3的层级关系
// 如果l1 > l2, 返回1，表示l1比l2高一个层级
// 如果l1 < l2, 返回-1，表示l1比l2低一个层级
// 如果l1 == l2, 返回0，表示l1和l2在同一层级
// Project > Module > Block > Point
func CompareLevel(l1, l2 string) int {
	levelMap := map[string]int{
		Project: 4,
		Module:  3,
		Block:   2,
		Point:   1,
	}

	l1Level := levelMap[l1]
	l2Level := levelMap[l2]

	switch {
	case l1Level > l2Level:
		return 1
	case l1Level < l2Level:
		return -1
	default:
		return 0
	}
}
