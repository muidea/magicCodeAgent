package service

import (
	"context"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	imkkc "github.com/muidea/magicCodeAgent/internal/modules/kernel/knowledgebase/pkg/common"
)

func (s *KnowledgeBase) CreateWarehouse(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.CreateWarehouseResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	param := &imkkc.CreateWarehouseParam{}
	if err := fn.ParseJSONBody(req, nil, param); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}
	if err := s.bizPtr.NewWarehouse(param.Name); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, err.Error())
		log.Errorf("new warehouse error: %s", err.Error())
		return
	}

	result.Error = nil
}

func (s *KnowledgeBase) DestroyWarehouse(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkkc.DestroyWarehouseResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(res, result)

	param := &imkkc.DestroyWarehouseParam{}
	if err := fn.ParseJSONBody(req, nil, param); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("parse body error: %s", err.Error())
		return
	}

	if err := s.bizPtr.DropWarehouse(param.Name); err != nil {
		result.Error = cd.NewError(cd.IllegalParam, "failed to parse request parameters")
		log.Errorf("drop warehouse error: %s", err.Error())
		return
	}

	result.Error = nil
}

func (s *KnowledgeBase) ListWarehouse(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	result := &imkkc.ListWarehouseResult{Result: cd.NewResult()}
	defer fn.PackageHTTPResponse(rw, result)
	valList, valErr := s.bizPtr.ListWarehouse()
	if valErr != nil {
		result.Error = cd.NewError(cd.Unexpected, valErr.Error())
		return
	}
	result.Values = valList
	result.Error = nil
}
