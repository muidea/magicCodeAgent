package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
)

func (s *Authority) filterAuthorityNamespaceLite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *fu.ContentFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.NamespaceLiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespaceList, _, namespaceErr := s.bizPtr.FilterAuthorityNamespaceLite(curSession.GetSessionInfo(), filter, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespaceList
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) FilterAuthorityNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)

	modelVal, modelOK := filter.Get("mode")
	if modelOK && modelVal == "1" {
		s.filterAuthorityNamespaceLite(ctx, res, req, filter)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.NamespaceStatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespaceList, namespaceTotal, namespaceErr := s.bizPtr.FilterAuthorityNamespace(curSession.GetSessionInfo(), filter, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespaceList
		result.Total = namespaceTotal
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) QueryAuthorityNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespacePtr, namespaceErr := s.bizPtr.QueryAuthorityNamespace(curSession.GetSessionInfo(), id, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		result.Namespace = namespacePtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) CreateAuthorityNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.NamespaceResult{}
	for {
		param := &cc.NamespaceParam{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespacePtr, namespaceErr := s.bizPtr.CreateAuthorityNamespace(curSession.GetSessionInfo(), param, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("新建Namespace:%s", namespacePtr.Name)
		s.writeLog(ctx, res, req, memo)

		result.Namespace = namespacePtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) UpdateAuthorityNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)

	result := &cc.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &cc.NamespaceParam{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespacePtr, namespaceErr := s.bizPtr.UpdateAuthorityNamespace(curSession.GetSessionInfo(), id, param, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("更新Namespace:%s", namespacePtr.Name)
		s.writeLog(ctx, res, req, memo)

		result.Namespace = namespacePtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Authority) DeleteAuthorityNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &cc.NamespaceResult{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		namespacePtr, namespaceErr := s.bizPtr.DeleteAuthorityNamespace(curSession.GetSessionInfo(), id, curNamespace)
		if namespaceErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = namespaceErr.Error()
			break
		}

		memo := fmt.Sprintf("删除Namespace:%s", namespacePtr.Name)
		s.writeLog(ctx, res, req, memo)

		result.Namespace = namespacePtr
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
