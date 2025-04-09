package biz

import (
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cClnt "github.com/muidea/magicCas/client"
	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) FilterAuthorityEndpoint(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.EndpointView, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterEndpoint(filter)
	return
}

func (s *Authority) FilterAuthorityEndpointLite(sessionInfo *session.SessionInfo, filter *fu.ContentFilter, namespace string) (ret []*cc.EndpointLite, total int64, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, total, err = clnt.FilterEndpointLite(filter)
	return
}

func (s *Authority) QueryAuthorityEndpoint(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EndpointView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.QueryEndpoint(id)
	return
}

func (s *Authority) CreateAuthorityEndpoint(sessionInfo *session.SessionInfo, ptr *cc.EndpointParam, namespace string) (ret *cc.EndpointView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.CreateEndpoint(ptr)
	return
}

func (s *Authority) UpdateAuthorityEndpoint(sessionInfo *session.SessionInfo, id int, ptr *cc.EndpointParam, namespace string) (ret *cc.EndpointView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	endpoint := ptr.ToEndpoint()
	endpoint.Namespace = namespace
	ret, err = clnt.UpdateEndpoint(id, ptr)
	return
}

func (s *Authority) DeleteAuthorityEndpoint(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EndpointView, err error) {
	clnt := cClnt.NewClient(s.casService)
	defer clnt.Release()

	clnt.BindSession(sessionInfo)
	clnt.AttachNameSpace(namespace)

	ret, err = clnt.DeleteEndpoint(id)
	return
}
