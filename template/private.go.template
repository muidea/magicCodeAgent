package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muidea/magicCas/client"
	casModel "github.com/muidea/magicCas/model"
	commonConst "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
)

func (s *Registry) queryPrivateGroup(currentSession session.Session) (ret []*casModel.PrivateGroupView, total int64, err error) {
	clnt := client.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, total, err = clnt.QueryAllPrivate()

	return
}

// QueryPrivateGroup query privateGroup
func (s *Registry) QueryPrivateGroup(res http.ResponseWriter, req *http.Request) {
	type queryResult struct {
		commonDef.Result
		PrivateGroups []*casModel.PrivateGroupView `json:"privates"`
	}

	result := &queryResult{}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		privateViewList, _, privateErr := s.queryPrivateGroup(curSession)
		if privateErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = fmt.Sprintf("查询权限组失败, %s", privateErr.Error())
			break
		}

		result.PrivateGroups = privateViewList
		break
	}
	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Registry) savePrivateGroup(currentSession session.Session, name, desc string, privateItems []*casModel.PrivateItem) (ret *casModel.PrivateGroupView, err error) {
	clnt := client.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	ret, err = clnt.SavePrivate(name, desc, privateItems)

	return
}

// SavePrivateGroup save privateGroup
func (s *Registry) SavePrivateGroup(res http.ResponseWriter, req *http.Request) {
	type saveParam struct {
		Name        string                  `json:"name"`
		Description string                  `json:"description"`
		Privates    []*casModel.PrivateItem `json:"privates"`
	}

	type saveResult struct {
		commonDef.Result
		PrivateGroup *casModel.PrivateGroupView `json:"private"`
	}

	result := &saveResult{}
	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		param := &saveParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "保存权限组失败, 非法参数"
			break
		}
		if param.Name == "" {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "保存权限组失败, 无效参数"
			break
		}

		privateView, privateErr := s.savePrivateGroup(curSession, param.Name, param.Description, param.Privates)
		if privateErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = fmt.Sprintf("保存权限组失败, %s", privateErr.Error())
			break
		}

		result.PrivateGroup = privateView

		memo := fmt.Sprintf("保存权限组%s", privateView.Name)
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

func (s *Registry) destoryPrivateGroup(currentSession session.Session, groupName string) (err error) {
	clnt := client.NewClient(s.casService)
	defer clnt.Release()

	sessionVal, ok := currentSession.GetOption(commonConst.SessionIdentity)
	if !ok {
		err = fmt.Errorf("无效权限")
		return
	}

	clnt.BindSession(sessionVal.(*commonConst.SessionInfo))
	err = clnt.DestroyPrivate(groupName)

	return
}

// DestoryPrivateGroup destory privateGroup
func (s *Registry) DestoryPrivateGroup(res http.ResponseWriter, req *http.Request) {
	result := &commonDef.Result{}

	curSession := s.sessionRegistry.GetSession(res, req)
	for {
		groupName := req.URL.Query().Get("groupName")
		if groupName == "" {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "删除权限组失败, 参数非法"
			break
		}

		err := s.destoryPrivateGroup(curSession, groupName)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = fmt.Sprintf("删除权限组失败, %s", err.Error())
			break
		}

		memo := fmt.Sprintf("删除权限组, Name:%s", groupName)
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
