package service

import (
	"context"
	"net/http"

	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
)

type Agent struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.Agent
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.Agent) *Agent {
	ptr := &Agent{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *Agent) RegisterRoute() {
	s.routeRegistry.SetApiVersion("")
	s.routeRegistry.AddHandler(common.ChatCompletion, engine.POST, s.ChatCompletion)
	s.routeRegistry.AddHandler(common.BetaCompletion, engine.POST, s.BetaCompletion)
	s.routeRegistry.AddHandler(common.ModelsList, engine.GET, s.ModelsList)
	s.routeRegistry.AddHandler(common.UserBalance, engine.GET, s.UserBalance)
}

func (s *Agent) ChatCompletion(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	paramPtr := &common.ChatCompletionRequest{}
	if err := fn.ParseJSONBody(req, nil, paramPtr); err != nil {
		log.Errorf("parse json body error:%s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if sse.IsSSE(req) || paramPtr.IsStream() {
		holder := sse.NewHolder(res, req)
		sseWriterContext := context.WithValue(ctx, common.SSEWriter{}, holder)
		taskFunc := func() error {
			return s.bizPtr.ChatCompletionStream("anonymous", sseWriterContext, paramPtr)
		}
		err := holder.Run(taskFunc)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		//res.WriteHeader(http.StatusOK)
		return
	}

	resultPtr, resultErr := s.bizPtr.ChatCompletion("anonymous", ctx, paramPtr)
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}

func (s *Agent) BetaCompletion(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	paramPtr := &common.BetaCompletionRequest{}
	if err := fn.ParseJSONBody(req, nil, paramPtr); err != nil {
		log.Errorf("parse json body error:%s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if sse.IsSSE(req) || paramPtr.IsStream() {
		holder := sse.NewHolder(res, req)
		sseWriterContext := context.WithValue(ctx, common.SSEWriter{}, holder)
		taskFunc := func() error {
			return s.bizPtr.BetaCompletionStream("anonymous", sseWriterContext, paramPtr)
		}
		err := holder.Run(taskFunc)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		//res.WriteHeader(http.StatusOK)
		return
	}

	resultPtr, resultErr := s.bizPtr.BetaCompletion("anonymous", ctx, paramPtr)
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}

func (s *Agent) ModelsList(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	resultPtr, resultErr := s.bizPtr.ModelsList()
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}

func (s *Agent) UserBalance(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	resultPtr, resultErr := s.bizPtr.UserBalance()
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}
