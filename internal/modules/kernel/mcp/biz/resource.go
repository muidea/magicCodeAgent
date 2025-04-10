package biz

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicToolkit/internal/config"
	mc "github.com/muidea/magicToolkit/internal/modules/kernel/monitor/pkg/common"
)

const mimeType = "application/json; charset=utf-8"

const (
	httpSchema = "http"
)

const (
	workspaceResource = "workspace"
	workspaceFile     = "file"
)

func (s *MCP) workspaceResourceURI() string {
	return fmt.Sprintf("%s://%s/%s", httpSchema, mcpServerName, workspaceResource)
}

func (s *MCP) workspaceResourceContentURI(name string) string {
	return fmt.Sprintf("%s://%s/%s/%s", httpSchema, mcpServerName, workspaceResource, name)
}

func (s *MCP) workspaceFileResourceURI() string {
	return fmt.Sprintf("%s://%s/%s/{id}/%s", httpSchema, mcpServerName, workspaceResource, workspaceFile)
}

// extractID extracts the ID from the URI(workspaceFileResourceURI)
func (s *MCP) extractID(uri string) (ret string, err error) {
	return
}

func (s *MCP) workspaceFileResourceContentURI(name string) string {
	return fmt.Sprintf("%s://%s/%s/%s/%s", httpSchema, mcpServerName, workspaceResource, name, workspaceFile)
}

func (s *MCP) resourceSet() {
	s.mcpServerPtr.AddResource(
		mcp.NewResource(
			s.workspaceResourceURI(),
			workspaceResource,
			mcp.WithMIMEType(mimeType),
		),
		s.workspaceResourceHandler,
	)

	s.mcpServerPtr.AddResourceTemplate(
		mcp.NewResourceTemplate(
			s.workspaceFileResourceURI(),
			workspaceFile,
			mcp.WithTemplateMIMEType(mimeType),
		),
		s.workspaceFileResourceHandler,
	)
}

func (s *MCP) workspaceResourceHandler(
	ctx context.Context,
	request mcp.ReadResourceRequest) (ret []mcp.ResourceContents, err error) {
	ev := event.NewEvent(mc.FilterWorkspace, s.ID(), mc.MonitorModule, nil, nil)
	resultVal, resultErr := s.SendEvent(ev).Get()
	if resultErr != nil {
		err = errors.New(resultErr.Error())
		log.Errorf("FilterWorkspace failed, s.SendEvent error:%s", err.Error())
		return
	}
	workspaceList := resultVal.([]*config.Workspace)
	for _, val := range workspaceList {
		ret = append(ret, mcp.TextResourceContents{
			URI:      s.workspaceResourceContentURI(val.Name),
			MIMEType: mimeType,
			Text:     s.ToText(val),
		})
	}
	return
}

func (s *MCP) workspaceFileResourceHandler(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) (ret []mcp.ResourceContents, err error) {
	if request.Params.URI == "" {
		err = errors.New("invalid request")
		return
	}
	idVal, idErr := s.extractID(request.Params.URI)
	if idErr != nil {
		err = idErr
		return
	}

	ev := event.NewEvent(mc.FilterWorkspaceFile, s.ID(), mc.MonitorModule, nil, idVal)
	resultVal, resultErr := s.SendEvent(ev).Get()
	if resultErr != nil {
		err = errors.New(resultErr.Error())
		log.Errorf("FilterWorkspaceFile failed, s.SendEvent error:%s", err.Error())
		return
	}
	ret = append(ret, mcp.TextResourceContents{
		URI:      s.workspaceFileResourceContentURI(idVal),
		MIMEType: mimeType,
		Text:     s.ToText(resultVal),
	})
	return
}
