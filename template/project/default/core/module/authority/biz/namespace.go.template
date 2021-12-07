package biz

import (
	casCommon "github.com/muidea/magicCas/common"
	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	casClient "github.com/muidea/magicCas/client"
)

func (s *Authority) FilterAuthorityNamespace(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.NamespaceView, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterNamespace(filter)
	return
}

func (s *Authority) FilterAuthorityNamespaceLite(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.NamespaceLite, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterNamespaceLite(filter)
	return
}

func (s *Authority) QueryAuthorityNamespace(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.NamespaceView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryNamespace(id)
	return
}

func (s *Authority) CreateAuthorityNamespace(sessionInfo *commonSession.SessionInfo, ptr *casCommon.NamespaceParam, namespace string) (ret *casCommon.NamespaceView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateNamespace(ptr)
	return
}

func (s *Authority) UpdateAuthorityNamespace(sessionInfo *commonSession.SessionInfo, id int, ptr *casCommon.NamespaceParam, namespace string) (ret *casCommon.NamespaceView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.UpdateNamespace(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityNamespace(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.NamespaceView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteNamespace(id)
	return
}
