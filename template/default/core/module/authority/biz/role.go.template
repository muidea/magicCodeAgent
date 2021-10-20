package biz

import (
	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	casClient "github.com/muidea/magicCas/client"
	casCommon "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityRole(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.RoleView, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterRole(filter)
	return
}

func (s *Authority) FilterAuthorityRoleLite(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.RoleLite, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterRoleLite(filter)
	return
}

func (s *Authority) QueryAuthorityRole(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.RoleView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryRole(id)
	return
}

func (s *Authority) CreateAuthorityRole(sessionInfo *commonSession.SessionInfo, ptr *casCommon.RoleParam, namespace string) (ret *casCommon.RoleView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateRole(ptr)
	return
}

func (s *Authority) UpdateAuthorityRole(sessionInfo *commonSession.SessionInfo, id int, ptr *casCommon.RoleParam, namespace string) (ret *casCommon.RoleView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateRole(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityRole(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.RoleView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteRole(id)
	return
}
