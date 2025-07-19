package service

import (
	"context"
	"errors"
	"net/http"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicEngine/sse"

	imkwc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

func (s *Workspace) createChatWorkspaceContentSSE(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	if !sse.IsSSE(req) {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	workspaceName := req.URL.Query().Get("workspace")
	holder := s.sseHolderRegistry.NewHolder(res, req)
	taskFunc := func() error {
		return s.bizPtr.ChatWorkspaceContentSSE(ctx, workspaceName, holder)
	}

	err := holder.EchoSSEID()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = holder.Run(taskFunc)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Workspace) submitChatWorkspaceContentMessage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &imkwc.ChatMessage{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &imkwc.ChatMessageResult{}
	defer fn.PackageHTTPResponse(res, result)
	holder := s.sseHolderRegistry.GetHolder(res, req)
	if holder == nil {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	taskFunc := func() error {
		errVal := s.bizPtr.ChatWorkspaceContentMessage(ctx, param)
		if errVal != nil {
			log.Errorf("chat workspace content message error: %s", errVal.Error())
			return errors.New("chat workspace content message error")
		}

		return nil
	}
	err = holder.Run(taskFunc)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	result.Error = nil
}

func (s *Workspace) TerminateChatWorkspaceContent(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &imkwc.ChatMessage{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &imkwc.ChatMessageResult{}
	defer fn.PackageHTTPResponse(res, result)
	holder := s.sseHolderRegistry.GetHolder(res, req)
	if holder == nil {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	taskFunc := func() error {
		errVal := s.bizPtr.ChatWorkspaceContentTerminate(ctx, param)
		if errVal != nil {
			log.Errorf("chat workspace content terminate error: %s", errVal.Error())
			return errors.New("chat workspace content terminate error")
		}

		return nil
	}
	err = holder.Run(taskFunc)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	result.Error = nil
}

func (s *Workspace) FilterChatWorkspaceContentHistory(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.ChatHistoryResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceName := req.URL.Query().Get("workspace")
	result.Values, result.Error = s.bizPtr.ChatWorkspaceContentHistory(workspaceName)
}

func (s *Workspace) DebugChatWorkspaceContent(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &imkwc.ChatMessage{}
	err := fn.ParseJSONBody(req, nil, param)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result := &imkwc.DebugWorkspaceResult{}
	defer fn.PackageHTTPResponse(res, result)
	holder := s.sseHolderRegistry.GetHolder(res, req)
	if holder == nil {
		log.Errorf("invalid sse request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	taskFunc := func() error {
		errVal := s.bizPtr.DebugWorkspaceContent(ctx, param)
		if errVal != nil {
			log.Errorf("chat workspace content debug error: %s", errVal.Error())
			return errors.New("chat workspace content debug error")
		}

		return nil
	}
	err = holder.Run(taskFunc)
	if err != nil {
		result.Error = cd.NewError(cd.Unexpected, err.Error())
		return
	}

	result.Error = nil
}

func (s *Workspace) FilterChatWorkspaceContentFlow(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &imkwc.FilterFlowResult{}
	defer fn.PackageHTTPResponse(res, result)

	workspaceName := req.URL.Query().Get("workspace")
	result.Values, result.Error = s.bizPtr.FilterWorkspaceFlow(workspaceName)
}
