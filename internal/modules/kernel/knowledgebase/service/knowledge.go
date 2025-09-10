package service

import (
	"context"
	"net/http"
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) uploadKnowledgeCallback(ctx context.Context, res http.ResponseWriter, req *http.Request, filePath string, err error) {
	result := &imkkc.ListKnowledgeResult{Result: *cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	name := filepath.Base(filePath)
	description := req.PostFormValue("description")
	warehouse := req.PostFormValue("warehouse")
	err = s.bizPtr.InsertKnowlege(warehouse, name, description, filePath)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
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

func (s *KnowledgeBase) DeleteKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.ListKnowledgeResult{Result: *cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	warehouse := req.URL.Query().Get("warehouse")
	filePath := req.URL.Query().Get("filePath")
	if warehouse == "" || filePath == "" {
		result.Error = cd.NewError(cd.IllegalParam, "failed to get request parameters")
		log.Errorf("missing request parameters")
		return
	}
	if err := s.bizPtr.DeleteKnowledge(warehouse, filePath); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, err.Error())
		log.Errorf("delete knowledge error: %s", err.Error())
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

func (s *KnowledgeBase) ListKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.ListKnowledgeResult{Result: *cd.NewResult()}
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

func (s *KnowledgeBase) GetKnowledgeContent(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.GetKnowledgeContentResult{Result: *cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)
	filePath := req.URL.Query().Get("filePath")
	if filePath == "" {
		result.Error = cd.NewError(cd.IllegalParam, "missing required parameter: filePath")
		log.Errorf("GetKnowledgeContent failed, missing filePath parameter")
		return
	}
	contentVal, contentErr := s.bizPtr.GetKnowledgeContent(filePath)
	if contentErr != nil {
		result.Error = cd.NewError(cd.Unexpected, contentErr.Error())
		log.Errorf("get knowledge content error: %s", contentErr.Error())
		return
	}

	result.Value = contentVal
	result.Error = nil
}

func (s *KnowledgeBase) SearchKnowledge(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.SearchKnowledgeResult{Result: *cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	param := &imkkc.SearchKnowledgeParam{}
	if err := fn.ParseJSONBody(req, nil, param); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	if param.Warehouse == "" || len(param.Text) == 0 {
		result.Error = cd.NewError(cd.IllegalParam, "warehouse and text can not be empty")
		log.Errorf("warehouse and text can not be empty")
		return
	}

	if param.Threshold <= 0.0001 {
		param.Threshold = 0.6
	}

	var chunkList []*imkkc.Chunk
	var chunkErr error
	if len(param.Text) > 1 {
		chunkList, chunkErr = s.bizPtr.SearchKnowledgeBulk(param.Warehouse, param.Text, param.TopK, param.Threshold)
	} else {
		chunkList, chunkErr = s.bizPtr.SearchKnowledge(param.Warehouse, param.Text[0], param.TopK, param.Threshold)
	}
	if chunkErr != nil {
		result.Error = cd.NewError(cd.Unexpected, "failed to search knowledge")
		log.Errorf("search knowledge error: %s", chunkErr.Error())
		return
	}
	result.Values = chunkList
	result.Error = nil
}
