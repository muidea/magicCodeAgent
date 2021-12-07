package biz

import (
	commonDef "github.com/muidea/magicCommon/def"
	commonSession "github.com/muidea/magicCommon/session"

	casClient "github.com/muidea/magicCas/client"
	casCommon "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityEndpoint(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.EndpointView, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterEndpoint(filter)
	return
}

func (s *Authority) FilterAuthorityEndpointLite(sessionInfo *commonSession.SessionInfo, filter *commonDef.Filter, namespace string) (ret []*casCommon.EndpointLite, total int64, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterEndpointLite(filter)
	return
}

func (s *Authority) QueryAuthorityEndpoint(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.EndpointView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryEndpoint(id)
	return
}

func (s *Authority) CreateAuthorityEndpoint(sessionInfo *commonSession.SessionInfo, ptr *casCommon.EndpointParam, namespace string) (ret *casCommon.EndpointView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateEndpoint(ptr)
	return
}

func (s *Authority) UpdateAuthorityEndpoint(sessionInfo *commonSession.SessionInfo, id int, ptr *casCommon.EndpointParam, namespace string) (ret *casCommon.EndpointView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	endpoint := ptr.ToEndpoint()
	endpoint.Namespace = namespace
	ret, err = clnt.UpdateEndpoint(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityEndpoint(sessionInfo *commonSession.SessionInfo, id int, namespace string) (ret *casCommon.EndpointView, err error) {
	clnt := casClient.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteEndpoint(id)
	return
}
