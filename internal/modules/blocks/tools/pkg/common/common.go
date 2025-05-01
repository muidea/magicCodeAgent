package common

import (
	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
)

const ToolsModule = "/internal/modules/blocks/tools"

type ToolParam map[string]any

type Tool interface {
	iipc.Prompts
	Execute(param ToolParam) error
}
