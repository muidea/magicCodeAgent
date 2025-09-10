package service

import (
	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/biz"
	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
)

type Workspace struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Workspace

	sseHolderRegistry *sse.HolderRegistry
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Workspace) *Workspace {
	ptr := &Workspace{
		routeRegistry:     routeRegistry,
		bizPtr:            bizPtr,
		sseHolderRegistry: sse.CreateHolderRegistry(),
	}
	return ptr
}

func (s *Workspace) RegisterRoute() {
	s.routeRegistry.SetApiVersion(common.ApiVersion)
	s.routeRegistry.AddHandler(imkwc.QueryWorkspace, engine.GET, s.queryWorkspace)
	s.routeRegistry.AddHandler(imkwc.CreateWorkspace, engine.POST, s.createWorkspace)
	s.routeRegistry.AddHandler(imkwc.DestroyWorkspace, engine.DELETE, s.destroyWorkspace)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspace, engine.GET, s.filterWorkspace)
	s.routeRegistry.AddHandler(imkwc.AddWorkspace, engine.POST, s.addWorkspace)
	s.routeRegistry.AddHandler(imkwc.UpdateWorkspace, engine.PUT, s.updateWorkspace)
	s.routeRegistry.AddHandler(imkwc.RemoveWorkspace, engine.DELETE, s.delWorkspace)
	s.routeRegistry.AddHandler(imkwc.EnableWorkspace, engine.PUT, s.enableWorkspace)
	s.routeRegistry.AddHandler(imkwc.DisableWorkspace, engine.PUT, s.disableWorkspace)
	s.routeRegistry.AddHandler(imkwc.RefreshWorkspace, engine.GET, s.refreshWorkspace)
	s.routeRegistry.AddHandler(imkwc.OpenWorkspaceExplorer, engine.GET, s.openWorkspace)
	s.routeRegistry.AddHandler(imkwc.SwitchWorkspaceExplorerStage, engine.PUT, s.switchWorkspaceStage)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceExplorerFile, engine.GET, s.filterFile)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceExplorerFileTree, engine.GET, s.filterFileTree)
	s.routeRegistry.AddHandler(imkwc.QueryWorkspaceExplorerFileContent, engine.GET, s.queryFileContent)
	s.routeRegistry.AddHandler(imkwc.CreateChatStreamSSE, engine.GET, s.createChatWorkspaceStreamSSE)
	s.routeRegistry.AddHandler(imkwc.SubmitChatStreamMessage, engine.POST, s.submitChatWorkspaceStreamMessage)
	s.routeRegistry.AddHandler(imkwc.TerminateChat, engine.POST, s.terminateChatWorkspaceContent)
	s.routeRegistry.AddHandler(imkwc.FlterChatHistory, engine.GET, s.filterChatWorkspaceContentHistory)
	s.routeRegistry.AddHandler(imkwc.DebugChatMessage, engine.POST, s.debugChatWorkspaceContent)
	s.routeRegistry.AddHandler(imkwc.FilterChatFlow, engine.GET, s.filterChatWorkspaceContentFlow)
	s.routeRegistry.AddHandler(imkwc.QuerySystemConfig, engine.GET, s.querySystemConfig)
	s.routeRegistry.AddHandler(imkwc.UpdateSystemConfig, engine.POST, s.updateSystemConfig)
}
