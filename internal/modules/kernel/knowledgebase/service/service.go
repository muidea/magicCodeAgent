package service

import (
	"context"
	"net/http"
	"path/filepath"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/biz"
	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

type KnowledgeBase struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.KnowledgeBase
	basePath      string
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.KnowledgeBase, basePath string) *KnowledgeBase {
	ptr := &KnowledgeBase{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
		basePath:      basePath,
	}
	return ptr
}

func (s *KnowledgeBase) RegisterRoute() {
	s.routeRegistry.AddHandler(imkkc.CreateKnowledgebaseWarehouse, engine.POST, s.CreateWarehouse)
	s.routeRegistry.AddHandler(imkkc.DestroyKnowledgebaseWarehouse, engine.DELETE, s.DestroyWarehouse)
	s.routeRegistry.AddHandler(imkkc.ListKnowledgebaseWarehouse, engine.GET, s.ListWarehouse)

	baseStaticPath := filepath.Join(s.basePath, imkkc.StaticPath)
	uploadRoute := engine.CreateUploadRoute(
		imkkc.UploadWorkspaceKnowledge,
		engine.POST,
		baseStaticPath,
		"",
		engine.DefaultMaxFileSize)
	s.routeRegistry.AddRoute(uploadRoute, s)
	s.routeRegistry.AddHandler(imkkc.ImportWorkspaceKnowledge, engine.POST, s.ImportKnowledge)
	s.routeRegistry.AddHandler(imkkc.DeleteWorkspaceKnowledge, engine.DELETE, s.DeleteKnowledge)
	s.routeRegistry.AddHandler(imkkc.ListWorkspaceKnowledge, engine.GET, s.ListKnowledge)
}

func (s *KnowledgeBase) MiddleWareHandle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	workspace := req.URL.Query().Get("workspace")
	if workspace == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	newContext := context.WithValue(ctx.Context(), engine.RelativePath{}, workspace)
	ctx.Update(newContext)
	ctx.Next()
}
