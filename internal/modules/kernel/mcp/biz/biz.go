package biz

import (
	"encoding/json"

	"github.com/mark3labs/mcp-go/server"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/mcp/pkg/common"
)

const mcpServerName = "mulife.vip/magicCodeAgent"
const mcpServerVersion = "1.0.0"

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
		mcpServerPtr: NewMCPServer(mcpServerName, mcpServerVersion),
	}
	return ptr
}

func (s *MCP) Initialize() (err *cd.Error) {
	s.resourceSet()

	return
}

func (s *MCP) GetMCPServer() *server.MCPServer {
	return s.mcpServerPtr
}

func (s *MCP) ToText(val any) (ret string) {
	byteVal, byteErr := json.MarshalIndent(val, "  ", "  ")
	if byteErr != nil {
		return
	}

	ret = string(byteVal)
	return
}
