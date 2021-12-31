package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) filterAuthorityEndpointLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *fu.ContentFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointList, _, endpointErr := s.bizPtr.FilterAuthorityEndpointLite(curSession.GetSessionInfo(), filter, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		result.Endpoint = endpointList
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) FilterAuthorityEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityEndpointLite(ctx, res, req, filter)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointList, endpointTotal, endpointErr := s.bizPtr.FilterAuthorityEndpoint(curSession.GetSessionInfo(), filter, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		roleList, _, roleErr := s.bizPtr.FilterAuthorityRoleLite(curSession.GetSessionInfo(), nil, curNamespace)
		if roleErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Endpoint = endpointList
		result.Role = roleList
		result.Total = endpointTotal
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) QueryAuthorityEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointPtr, endpointErr := s.bizPtr.QueryAuthorityEndpoint(curSession.GetSessionInfo(), id, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		result.Endpoint = endpointPtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) CreateAuthorityEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointResult{}
	for {
		param := &cc.EndpointParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointPtr, endpointErr := s.bizPtr.CreateAuthorityEndpoint(curSession.GetSessionInfo(), param, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(ctx, res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) UpdateAuthorityEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &cc.EndpointParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointPtr, endpointErr := s.bizPtr.UpdateAuthorityEndpoint(curSession.GetSessionInfo(), id, param, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(ctx, res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) DeleteAuthorityEndpoint(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		endpointPtr, endpointErr := s.bizPtr.DeleteAuthorityEndpoint(curSession.GetSessionInfo(), id, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(ctx, res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
