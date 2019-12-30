package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	casClient "github.com/muidea/magicCas/client"
	casModel "github.com/muidea/magicCas/model"
	commonConst "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
)

func (s *Registry) queryAllAccount(currentSession session.Session) (ret []*casModel.AccountView, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, total, err = clnt.QueryAllAccount(nil)

	return
}

// QueryAllAccount query member
func (s *Registry) QueryAllAccount(res http.ResponseWriter, req *http.Request) {
	type queryResult struct {
		commonDef.Result
		Total    int64                   `json:"total"`
		Accounts []*casModel.AccountView `json:"accounts"`
	}

	result := &queryResult{Accounts: []*casModel.AccountView{}}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		accountList, accountTotal, accountErr := s.queryAllAccount(curSession)
		if accountErr != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = accountErr.Error()
			break
		}

		result.Total = accountTotal
		result.Accounts = accountList
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

func (s *Registry) deleteAccount(currentSession session.Session, id int) (ret *casModel.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("删除账号失败，无效权限")
		return
	}
	authVal, ok := currentSession.GetOption(commonConst.AuthAccount)
	if !ok {
		err = fmt.Errorf("删除账号失败，无效权限")
		return
	}
	accountView := authVal.(*casModel.AccountView)
	if accountView.ID == id {
		err = fmt.Errorf("删除账号失败，禁止删除当前登录账号")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, err = clnt.DestroyAccount(id)
	return
}

// DeleteAccount delete account
func (s *Registry) DeleteAccount(res http.ResponseWriter, req *http.Request) {
	result := &commonDef.Result{}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("delete account failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		accountPtr, accountErr := s.deleteAccount(curSession, id)
		if accountErr != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = accountErr.Error()
			break
		}

		result.ErrorCode = commonDef.Success

		memo := fmt.Sprintf("删除账号%s", accountPtr.Account)
		s.writelog(res, req, memo)
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Registry) queryAccount(currentSession session.Session, id int) (ret *casModel.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, _, err = clnt.QueryAccount(id)
	return
}

// QueryAccount query account
func (s *Registry) QueryAccount(res http.ResponseWriter, req *http.Request) {
	type queryResult struct {
		commonDef.Result
		Account *casModel.AccountView `json:"account"`
	}

	result := &queryResult{}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("update account failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		accountView, accountErr := s.queryAccount(curSession, id)
		if accountErr != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "参数非法"
		}

		result.Account = accountView
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

func (s *Registry) updateAccount(currentSession session.Session, account *casModel.AccountView) (ret *casModel.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, err = clnt.UpdateAccount(account)
	return
}

// UpdateAccount update account
func (s *Registry) UpdateAccount(res http.ResponseWriter, req *http.Request) {
	param := map[string]int{}
	result := &commonDef.Result{}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("update account failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		accountView, accountErr := s.queryAccount(curSession, id)
		if accountErr != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "参数非法"
			break
		}

		err = net.ParseJSONBody(req, &param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "参数非法"
			break
		}
		status, ok := param["status"]
		if ok {
			accountView.Status = casModel.GetStatus(status)
		}
		privateGroup, ok := param["privateGroup"]
		if ok {
			accountView.PrivateGroup = &casModel.PrivateGroupLite{ID: privateGroup}
		}

		_, err = s.updateAccount(curSession, accountView)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "更新信息出错,更新账号信息失败"
			break
		}

		result.ErrorCode = commonDef.Success

		memo := fmt.Sprintf("更新账号%s", accountView.Account)
		s.writelog(res, req, memo)
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
