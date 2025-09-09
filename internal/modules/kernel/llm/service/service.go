package service

import (
	"context"
	"net/http"

	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	engine "github.com/muidea/magicEngine/http"
	"github.com/muidea/magicEngine/sse"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/biz"
	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
)

type LLM struct {
	routeRegistry engine.RouteRegistry
	bizPtr        *biz.LLM
}

func New(routeRegistry engine.RouteRegistry, bizPtr *biz.LLM) *LLM {
	ptr := &LLM{
		routeRegistry: routeRegistry,
		bizPtr:        bizPtr,
	}
	return ptr
}

func (s *LLM) RegisterRoute() {
	s.routeRegistry.SetApiVersion("")
	s.routeRegistry.AddHandler(imklc.ChatCompletionStream, engine.POST, s.ChatCompletionStream)
	s.routeRegistry.AddHandler(imklc.BetaCompletion, engine.POST, s.BetaCompletion)
	s.routeRegistry.AddHandler(imklc.ModelsList, engine.GET, s.ModelsList)
	s.routeRegistry.AddHandler(imklc.UserBalance, engine.GET, s.UserBalance)
}

func (s *LLM) ChatCompletionStream(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	paramPtr := &imklc.ChatCompletionRequest{}
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

func (s *LLM) BetaCompletion(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	paramPtr := &imklc.BetaCompletionRequest{}
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

func (s *LLM) ModelsList(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	resultPtr, resultErr := s.bizPtr.ModelsList()
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}

func (s *LLM) UserBalance(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	resultPtr, resultErr := s.bizPtr.UserBalance()
	if resultErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	fn.PackageHTTPResponse(res, resultPtr)
}
