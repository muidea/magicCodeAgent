package biz

import (
	"github.com/mark3labs/mcp-go/server"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/mcp/pkg/common"
)

type MCP struct {
	biz.Base

	mcpServerPtr *server.MCPServer
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *MCP {
	ptr := &MCP{
		Base:         biz.New(common.MCPModule, eventHub, backgroundRoutine),
		mcpServerPtr: NewMCPServer(),
	}
	return ptr
}

func (s *MCP) Initialize() (err *cd.Error) {
	return
}

func (s *MCP) GetMCPServer() *server.MCPServer {
	return s.mcpServerPtr
}
