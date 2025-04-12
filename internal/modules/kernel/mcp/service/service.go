package service

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/mcp/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/mcp/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"
)

type MCP struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.MCP
	sseServer     *server.SSEServer
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.MCP) *MCP {
	return &MCP{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
		sseServer: server.NewSSEServer(
			bizPtr.GetMCPServer(),
			server.WithBasePath(tc.ApiVersion),
			server.WithSSEEndpoint(common.MCPSSE),
			server.WithMessageEndpoint(common.MCPMessage),
		),
	}
}

func (s *MCP) RegisterRoute() {
	s.routeRegistry.AddHandler(common.MCPSSE, engine.GET, s.sseHandler)
	s.routeRegistry.AddHandler(common.MCPMessage, engine.POST, s.messageHandler)
}

func (s *MCP) sseHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.sseServer.ServeHTTP(w, r)
}

func (s *MCP) messageHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.sseServer.ServeHTTP(w, r)
}
