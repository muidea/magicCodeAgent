package service

import (
	"encoding/json"
	"fmt"
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	"net/http"
)

func (s *Authority) filterAuthorityNamespaceLite(res http.ResponseWriter, req *http.Request, filter *commonDef.Filter) {
	result := &casCommon.NamespaceLiteListResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespaceList, _, namespaceErr := s.bizPtr.FilterAuthorityNamespaceLite(curSession.GetSessionInfo(), filter, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespaceList
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

func (s *Authority) FilterAuthorityNamespace(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityNamespaceLite(res, req, filter)
		return
	}

	result := &casCommon.NamespaceStatisticResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespaceList, namespaceTotal, namespaceErr := s.bizPtr.FilterAuthorityNamespace(curSession.GetSessionInfo(), filter, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespaceList
		result.Total = namespaceTotal
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

func (s *Authority) QueryAuthorityNamespace(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespacePtr, namespaceErr := s.bizPtr.QueryAuthorityNamespace(curSession.GetSessionInfo(), id, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespacePtr
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

func (s *Authority) CreateAuthorityNamespace(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.NamespaceResult{}
	for {
		param := &casCommon.NamespaceParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespacePtr, namespaceErr := s.bizPtr.CreateAuthorityNamespace(curSession.GetSessionInfo(), param, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Namespace:%s", namespacePtr.Name)
		s.writeLog(res, req, memo)

		result.Namespace = namespacePtr
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

func (s *Authority) UpdateAuthorityNamespace(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		param := &casCommon.NamespaceParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespacePtr, namespaceErr := s.bizPtr.UpdateAuthorityNamespace(curSession.GetSessionInfo(), id, param, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Namespace:%s", namespacePtr.Name)
		s.writeLog(res, req, memo)

		result.Namespace = namespacePtr
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

func (s *Authority) DeleteAuthorityNamespace(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		namespacePtr, namespaceErr := s.bizPtr.DeleteAuthorityNamespace(curSession.GetSessionInfo(), id, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Namespace:%s", namespacePtr.Name)
		s.writeLog(res, req, memo)

		result.Namespace = namespacePtr
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
