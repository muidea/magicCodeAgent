package service

import (
	"encoding/json"
	"fmt"
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	"net/http"
)

func (s *Authority) filterAuthorityEndpointLite(res http.ResponseWriter, req *http.Request, filter *commonDef.Filter) {
	result := &casCommon.EndpointLiteListResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointList, _, endpointErr := s.bizPtr.FilterAuthorityEndpointLite(curSession.GetSessionInfo(), filter, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		result.Endpoint = endpointList
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) FilterAuthorityEndpoint(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityEndpointLite(res, req, filter)
		return
	}

	result := &casCommon.EndpointStatisticResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointList, endpointTotal, endpointErr := s.bizPtr.FilterAuthorityEndpoint(curSession.GetSessionInfo(), filter, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		roleList, _, roleErr := s.bizPtr.FilterAuthorityRoleLite(curSession.GetSessionInfo(), nil, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Endpoint = endpointList
		result.Role = roleList
		result.Total = endpointTotal
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) QueryAuthorityEndpoint(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointPtr, endpointErr := s.bizPtr.QueryAuthorityEndpoint(curSession.GetSessionInfo(), id, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		result.Endpoint = endpointPtr
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) CreateAuthorityEndpoint(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.EndpointResult{}
	for {
		param := &casCommon.EndpointParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointPtr, endpointErr := s.bizPtr.CreateAuthorityEndpoint(curSession.GetSessionInfo(), param, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) UpdateAuthorityEndpoint(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		param := &casCommon.EndpointParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointPtr, endpointErr := s.bizPtr.UpdateAuthorityEndpoint(curSession.GetSessionInfo(), id, param, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) DeleteAuthorityEndpoint(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.EndpointResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		endpointPtr, endpointErr := s.bizPtr.DeleteAuthorityEndpoint(curSession.GetSessionInfo(), id, curNamespace)
		if endpointErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = endpointErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Endpoint:%s", endpointPtr.Endpoint)
		s.writeLog(res, req, memo)

		result.Endpoint = endpointPtr
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
