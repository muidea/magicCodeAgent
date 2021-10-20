package service

import (
	"encoding/json"
	"fmt"
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	"net/http"
)

func (s *Authority) filterAuthorityAccountLite(res http.ResponseWriter, req *http.Request, filter *commonDef.Filter) {
	result := &casCommon.AccountLiteListResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		filter.Remove("mode")

		accountList, accountTotal, accountErr := s.bizPtr.FilterAuthorityAccountLite(curSession.GetSessionInfo(), filter, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		result.Account = accountList
		result.Total = accountTotal
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

func (s *Authority) FilterAuthorityAccount(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityAccountLite(res, req, filter)
		return
	}

	result := &casCommon.AccountStatisticResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		accountList, accountTotal, accountErr := s.bizPtr.FilterAuthorityAccount(curSession.GetSessionInfo(), filter, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		roleList, _, roleErr := s.bizPtr.FilterAuthorityRoleLite(curSession.GetSessionInfo(), nil, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Account = accountList
		result.Role = roleList
		result.Total = accountTotal
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

func (s *Authority) QueryAuthorityAccount(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.AccountResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		accountPtr, accountErr := s.bizPtr.QueryAuthorityAccount(curSession.GetSessionInfo(), id, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		result.Account = accountPtr
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

func (s *Authority) CreateAuthorityAccount(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.AccountResult{}
	for {
		param := &casCommon.AccountParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		accountPtr, accountErr := s.bizPtr.CreateAuthorityAccount(curSession.GetSessionInfo(), param, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Account:%s", accountPtr.Account)
		s.writeLog(res, req, memo)

		result.Account = accountPtr
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

func (s *Authority) UpdateAuthorityAccount(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.AccountResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		param := &casCommon.AccountParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		accountPtr, accountErr := s.bizPtr.UpdateAuthorityAccount(curSession.GetSessionInfo(), id, param, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Account:%s", accountPtr.Account)
		s.writeLog(res, req, memo)

		result.Account = accountPtr
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

func (s *Authority) DeleteAuthorityAccount(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.AccountResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		accountPtr, accountErr := s.bizPtr.DeleteAuthorityAccount(curSession.GetSessionInfo(), id, curNamespace)
		if accountErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = accountErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Account:%s", accountPtr.Account)
		s.writeLog(res, req, memo)

		result.Account = accountPtr
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
