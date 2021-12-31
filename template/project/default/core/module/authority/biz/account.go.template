package biz

import (
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cClnt "github.com/muidea/magicCas/client"
	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityAccount(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.AccountView, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterAccount(filter)
	return
}

func (s *Authority) FilterAuthorityAccountLite(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.AccountLite, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterAccountLite(filter)
	return
}

func (s *Authority) QueryAuthorityAccount(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.AccountView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryAccount(id)
	return
}

func (s *Authority) CreateAuthorityAccount(sessionInfo *session.SessionInfo, ptr *cc.AccountParam, namespace string) (ret *cc.AccountView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	param := &cc.AccountParam{Account: ptr.Account, Password: ptr.Password, EMail: ptr.EMail, Description: ptr.Description}
	ret, err = clnt.CreateAccount(param)
	return
}

func (s *Authority) UpdateAuthorityAccount(sessionInfo *session.SessionInfo, id int, ptr *cc.AccountParam, namespace string) (ret *cc.AccountView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateAccount(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityAccount(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.AccountView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteAccount(id)
	return
}
