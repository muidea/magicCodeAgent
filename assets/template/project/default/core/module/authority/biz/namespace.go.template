package biz

import (
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cClnt "github.com/muidea/magicCas/client"
	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityNamespace(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.NamespaceView, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterNamespace(filter)
	return
}

func (s *Authority) FilterAuthorityNamespaceLite(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.NamespaceLite, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterNamespaceLite(filter)
	return
}

func (s *Authority) QueryAuthorityNamespace(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.NamespaceView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryNamespace(id)
	return
}

func (s *Authority) CreateAuthorityNamespace(sessionInfo *session.SessionInfo, ptr *cc.NamespaceParam, namespace string) (ret *cc.NamespaceView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateNamespace(ptr)
	return
}

func (s *Authority) UpdateAuthorityNamespace(sessionInfo *session.SessionInfo, id int, ptr *cc.NamespaceParam, namespace string) (ret *cc.NamespaceView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateNamespace(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityNamespace(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.NamespaceView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteNamespace(id)
	return
}
