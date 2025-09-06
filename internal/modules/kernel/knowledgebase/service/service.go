package service

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
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
		engine.DefaultMaxFileSize, s.uploadKnowledgeCallback)
	s.routeRegistry.AddRoute(uploadRoute, s)
	s.routeRegistry.AddHandler(imkkc.ImportWorkspaceKnowledge, engine.POST, s.ImportKnowledge)
	s.routeRegistry.AddHandler(imkkc.DeleteWorkspaceKnowledge, engine.DELETE, s.DeleteKnowledge)
	s.routeRegistry.AddHandler(imkkc.ListWorkspaceKnowledge, engine.GET, s.ListKnowledge)
}

func (s *KnowledgeBase) MiddleWareHandle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(engine.DefaultMaxFileSize)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	warehouse := req.PostFormValue("warehouse")
	if warehouse == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !s.bizPtr.ExistWarehouse(warehouse) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("warehouse:%s", warehouse)
	newContext := context.WithValue(ctx.Context(), engine.RelativePath{}, warehouse)
	ctx.Update(newContext)
	ctx.Next()
}
