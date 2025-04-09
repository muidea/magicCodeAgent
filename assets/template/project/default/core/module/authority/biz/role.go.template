package biz

import (
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cClnt "github.com/muidea/magicCas/client"
	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityRole(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.RoleView, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterRole(filter)
	return
}

func (s *Authority) FilterAuthorityRoleLite(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.RoleLite, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterRoleLite(filter)
	return
}

func (s *Authority) QueryAuthorityRole(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.RoleView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryRole(id)
	return
}

func (s *Authority) CreateAuthorityRole(sessionInfo *session.SessionInfo, ptr *cc.RoleParam, namespace string) (ret *cc.RoleView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateRole(ptr)
	return
}

func (s *Authority) UpdateAuthorityRole(sessionInfo *session.SessionInfo, id int, ptr *cc.RoleParam, namespace string) (ret *cc.RoleView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateRole(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityRole(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.RoleView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteRole(id)
	return
}
