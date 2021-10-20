package service

import (
	"encoding/json"
	"fmt"
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	"net/http"
)

func (s *Authority) filterAuthorityRoleLite(res http.ResponseWriter, req *http.Request, filter *commonDef.Filter) {
	result := &casCommon.RoleLiteListResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		roleList, _, roleErr := s.bizPtr.FilterAuthorityRoleLite(curSession.GetSessionInfo(), filter, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Role = roleList
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

func (s *Authority) FilterAuthorityRole(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityRoleLite(res, req, filter)
		return
	}

	result := &casCommon.RoleStatisticResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		roleList, roleTotal, roleErr := s.bizPtr.FilterAuthorityRole(curSession.GetSessionInfo(), filter, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Role = roleList
		result.Total = roleTotal
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

func (s *Authority) QueryAuthorityRole(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.RoleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		rolePtr, roleErr := s.bizPtr.QueryAuthorityRole(curSession.GetSessionInfo(), id, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		result.Role = rolePtr
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

func (s *Authority) CreateAuthorityRole(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.RoleResult{}
	for {
		param := &casCommon.RoleParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		rolePtr, roleErr := s.bizPtr.CreateAuthorityRole(curSession.GetSessionInfo(), param, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Role:%s", rolePtr.Name)
		s.writeLog(res, req, memo)

		result.Role = rolePtr
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

func (s *Authority) UpdateAuthorityRole(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.RoleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		param := &casCommon.RoleParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		rolePtr, roleErr := s.bizPtr.UpdateAuthorityRole(curSession.GetSessionInfo(), id, param, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Role:%s", rolePtr.Name)
		s.writeLog(res, req, memo)

		result.Role = rolePtr
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

func (s *Authority) DeleteAuthorityRole(res http.ResponseWriter, req *http.Request) {
	result := &casCommon.RoleResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		rolePtr, roleErr := s.bizPtr.DeleteAuthorityRole(curSession.GetSessionInfo(), id, curNamespace)
		if roleErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = roleErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Role:%s", rolePtr.Name)
		s.writeLog(res, req, memo)

		result.Role = rolePtr
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
