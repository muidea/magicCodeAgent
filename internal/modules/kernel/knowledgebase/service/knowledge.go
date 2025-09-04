package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) ImportKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := cd.NewResult()
	defer fn.PackageHTTPResponse(res, result)

	param := &imkkc.ImportKnowledgeParam{}
	if err := fn.ParseJSONBody(req, nil, param); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	err := s.bizPtr.ImportKnowledge(param.FilePath, param.Warehouse)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		log.Errorf("import knowledge error: %s", err.Error())
		return
	}
	result.Error = nil
}

func (s *KnowledgeBase) DeleteKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := cd.NewResult()
	defer fn.PackageHTTPResponse(res, result)

	param := &imkkc.DeleteKnowledgeParam{}
	if err := fn.ParseJSONBody(req, nil, param); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	if err := s.bizPtr.DeleteKnowledge(param.Warehouse, param.FilePath); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, err.Error())
		log.Errorf("delete knowledge error: %s", err.Error())
		return
	}
	result.Error = nil
}

func (s *KnowledgeBase) ListKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.ListKnowledgeResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	warehouse := req.URL.Query().Get("warehouse")
	if warehouse == "" {
		result.Error = cd.NewError(cd.IllegalParam, "missing required parameter: warehouse")
		log.Errorf("ListKnowledge failed, missing warehouse parameter")
		return
	}

	valList, valErr := s.bizPtr.ListKnowledge(warehouse)
	if valErr != nil {
		result.Error = cd.NewError(cd.Unexpected, valErr.Error())
		log.Errorf("list knowledge error: %s", valErr.Error())
		return
	}

	result.Values = valList
	result.Error = nil
}
