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
	Recorder  = "recorder"
)

const (
	CodeStage   = "code"
	DesignStage = "design"
)
