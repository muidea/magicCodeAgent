package service

import (
	"context"
	"net/http"

	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/biz"
	imbtmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
)

type Templatemanager struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Templatemanager
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Templatemanager) *Templatemanager {
	ptr := &Templatemanager{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Templatemanager) RegisterRoute() {
	s.routeRegistry.AddHandler(imbtmc.ListTemplate, engine.GET, s.listTemplate)
	s.routeRegistry.AddHandler(imbtmc.LoadTemplate, engine.GET, s.loadTemplate)
}

func (s *Templatemanager) listTemplate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imbtmc.ListTemplateResult{}
	defer fn.PackageHTTPResponse(res, result)

	result.Values = s.bizPtr.ListTemplate()
}

func (s *Templatemanager) loadTemplate(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imbtmc.LoadTemplateResult{}
	defer fn.PackageHTTPResponse(res, result)

	templateName := req.URL.Query().Get("name")
	result.Value, result.Error = s.bizPtr.LoadTemplate(templateName)
}
