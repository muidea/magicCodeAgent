package service

import (
	"context"
	"net/http"

	engine "github.com/muidea/magicEngine/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Tracer struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Tracer
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Tracer) *Tracer {
	ptr := &Tracer{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Tracer) RegisterRoute() {
	s.routeRegistry.AddHandler(common.AppendMessage, engine.POST, s.AppendMessage)
	s.routeRegistry.AddHandler(common.FilterMessage, engine.GET, s.FilterMessage)
}

func (s *Tracer) AppendMessage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := cd.NewResult()
	defer fn.PackageHTTPResponse(res, result)

	msg := &models.MessageInfo{}
	err := fn.ParseJSONBody(req, nil, msg)
	if err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "invalid append message parameters")
		return
	}
	s.bizPtr.AppendMessage(msg)
	result.Error = nil
}

func (s *Tracer) FilterMessage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.MessageInfoList{}
	defer fn.PackageHTTPResponse(res, result)

	filter := &common.MessageFilter{}
	filter.Decode(req)

	result.Values, result.Error = s.bizPtr.FilterMessage(filter)
}
