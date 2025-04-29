package service

import (
	"context"
	"net/http"

	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/biz"
	imkac "github.com/muidea/magicCodeAgent/internal/modules/kernel/agent/pkg/common"
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
	s.routeRegistry.AddHandler(imkac.ChatCompletion, engine.POST, s.ChatCompletion)
	s.routeRegistry.AddHandler(imkac.BetaCompletion, engine.POST, s.BetaCompletion)
	s.routeRegistry.AddHandler(imkac.ModelsList, engine.GET, s.ModelsList)
	s.routeRegistry.AddHandler(imkac.UserBalance, engine.GET, s.UserBalance)
}

func (s *Agent) ChatCompletion(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	paramPtr := &imkac.ChatCompletionStreamRequest{}
	if err := fn.ParseJSONBody(req, nil, paramPtr); err != nil {
		log.Errorf("parse json body error:%s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if sse.IsSSE(req) || paramPtr.IsStream() {
		holder := sse.NewHolder(res, req)
		sseWriterContext := context.WithValue(ctx, sse.StreamSink{}, holder)
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
	paramPtr := &imkac.BetaCompletionRequest{}
	if err := fn.ParseJSONBody(req, nil, paramPtr); err != nil {
		log.Errorf("parse json body error:%s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if sse.IsSSE(req) || paramPtr.IsStream() {
		holder := sse.NewHolder(res, req)
		sseWriterContext := context.WithValue(ctx, sse.StreamSink{}, holder)
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
