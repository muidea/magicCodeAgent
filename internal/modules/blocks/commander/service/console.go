package service

import (
	"context"
	"net/http"

	fn "github.com/muidea/magicCommon/foundation/net"
)

func (s *Commander) ConsoleLogs(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	byteVal, byteErr := fn.GetHTTPRequestBody(req)
	if byteErr != nil {
		res.WriteHeader(http.StatusOK)
		return
	}

	s.bizPtr.AppendConsoleLogs("", string(byteVal))
	res.WriteHeader(http.StatusOK)
}
