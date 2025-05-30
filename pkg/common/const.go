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
