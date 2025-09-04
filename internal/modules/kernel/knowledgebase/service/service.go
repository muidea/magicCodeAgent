package service

import (
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

	baseStaticPath := filepath.Join(s.basePath, imkkc.StaticPath)
	uploadRoute := engine.CreateUploadRoute(
		imkkc.UploadWorkspaceKnowledge,
		engine.POST,
		baseStaticPath,
		"",
		engine.DefaultMaxFileSize)
	s.routeRegistry.AddRoute(uploadRoute)
	s.routeRegistry.AddHandler(imkkc.ImportWorkspaceKnowledge, engine.POST, s.ImportKnowledge)
}
