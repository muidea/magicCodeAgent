package biz

import (
	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	casClient "github.com/muidea/magicCas/client"
	casCommon "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityAccount(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.AccountView, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterAccount(filter)
	return
}

func (s *Authority) FilterAuthorityAccountLite(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.AccountLite, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterAccountLite(filter)
	return
}

func (s *Authority) QueryAuthorityAccount(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryAccount(id)
	return
}

func (s *Authority) CreateAuthorityAccount(sessionInfo *commonSession.SessionInfo, ptr *casCommon.AccountParam, namespace string) (ret *casCommon.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	param := &casCommon.AccountParam{Account: ptr.Account, Password: ptr.Password, EMail: ptr.EMail, Description: ptr.Description}
	ret, err = clnt.CreateAccount(param)
	return
}

func (s *Authority) UpdateAuthorityAccount(sessionInfo *commonSession.SessionInfo, id int, ptr *casCommon.AccountParam, namespace string) (ret *casCommon.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateAccount(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityAccount(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.AccountView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteAccount(id)
	return
}
