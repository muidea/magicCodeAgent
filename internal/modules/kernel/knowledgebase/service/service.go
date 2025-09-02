package service

import (
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/biz"
	engine "github.com/muidea/magicEngine/http"
)

type KnowledgeBase struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.KnowledgeBase
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.KnowledgeBase) *KnowledgeBase {
	ptr := &KnowledgeBase{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *KnowledgeBase) RegisterRoute() {
}
