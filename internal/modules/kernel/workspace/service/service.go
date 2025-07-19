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
	s.routeRegistry.AddHandler(imkwc.UpdateWorkspaceStage, engine.PUT, s.updateWorkspaceStage)
	s.routeRegistry.AddHandler(imkwc.RefreshWorkspace, engine.GET, s.refreshWorkspace)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceFile, engine.GET, s.filterFile)
	s.routeRegistry.AddHandler(imkwc.FilterWorkspaceFileTree, engine.GET, s.filterFileTree)
	s.routeRegistry.AddHandler(imkwc.QueryWorkspaceFileContent, engine.GET, s.queryFileContent)
	s.routeRegistry.AddHandler(imkwc.CreateChatWorkspaceContentSSE, engine.GET, s.createChatWorkspaceContentSSE)
	s.routeRegistry.AddHandler(imkwc.SubmitChatWorkspaceContentMessage, engine.POST, s.submitChatWorkspaceContentMessage)
	s.routeRegistry.AddHandler(imkwc.TerminateChatWorkspaceContent, engine.POST, s.TerminateChatWorkspaceContent)
	s.routeRegistry.AddHandler(imkwc.FlterChatWorkspaceContentHistory, engine.GET, s.FilterChatWorkspaceContentHistory)
	s.routeRegistry.AddHandler(imkwc.DebugChatWorkspaceContent, engine.POST, s.DebugChatWorkspaceContent)
	s.routeRegistry.AddHandler(imkwc.FilterChatWorkspaceContentFlow, engine.GET, s.FilterChatWorkspaceContentFlow)
	s.routeRegistry.AddHandler(imkwc.QuerySystemConfig, engine.GET, s.querySystemConfig)
	s.routeRegistry.AddHandler(imkwc.UpdateSystemConfig, engine.POST, s.updateSystemConfig)
}
