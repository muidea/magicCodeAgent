package biz

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strings"

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
	urlVal, urlErr := url.Parse(mcpServerName)
	if urlErr != nil {
		return ""
	}
	urlVal.Scheme = httpSchema
	urlVal.Path = path.Join(urlVal.Path, workspaceResource)

	return urlVal.String()
}

func (s *MCP) workspaceResourceContentURI(name string) string {
	urlVal, urlErr := url.Parse(mcpServerName)
	if urlErr != nil {
		return ""
	}
	urlVal.Scheme = httpSchema
	urlVal.Path = path.Join(urlVal.Path, workspaceResource, name)

	return urlVal.String()
}

func (s *MCP) workspaceFileResourceURI() string {
	urlVal, urlErr := url.Parse(mcpServerName)
	if urlErr != nil {
		return ""
	}

	urlVal.Scheme = httpSchema
	urlVal.Path = path.Join(urlVal.Path, workspaceResource, "{id}", workspaceFile)

	return urlVal.String()
}

// extractIDFromFileResourceContentURI extracts the ID from the URI(workspaceFileResourceURI)
// 例如"http://mulife.vip/magicCodeAgent/workspace/test/file"，取出test
// "magicCodeAgent/workspace/test/file", 取出test
func (s *MCP) extractIDFromFileResourceContentURI(uri string) (ret string, err error) {
	urlVal, urlErr := url.Parse(uri)
	if urlErr != nil {
		log.Errorf("parse url failed, urlErr:%s", urlErr.Error())
		err = urlErr
		return
	}

	parts := strings.Split(urlVal.Path, "/")
	if len(parts) != 5 {
		log.Errorf("invalid url, url:%s", urlVal.Path)
		err = errors.New("invalid url")
		return
	}
	if parts[4] != "file" || parts[2] != "workspace" {
		log.Errorf("invalid url, url:%s", urlVal.Path)
		err = errors.New("invalid url")
		return
	}

	ret = parts[3]
	return
}

func (s *MCP) workspaceFileResourceContentURI(name string) string {
	urlVal, urlErr := url.Parse(mcpServerName)
	if urlErr != nil {
		log.Errorf("parse url failed, urlErr:%s", urlErr.Error())
		return ""
	}
	urlVal.Scheme = httpSchema
	urlVal.Path = path.Join(urlVal.Path, workspaceResource, name, workspaceFile)

	return urlVal.String()
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
	idVal, idErr := s.extractIDFromFileResourceContentURI(request.Params.URI)
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
