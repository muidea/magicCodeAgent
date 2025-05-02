package biz

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/internal/config"
	mc "github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

const mimeType = "application/json; charset=utf-8"

const (
	httpSchema = "http"
)

const (
	workspaceResource = "workspace"
	workspaceFile     = "file"

	// 工作区清单
	workspaceResourceName = "Workspace_list"
	// 工作区清单，查看当前设置的工作区清单
	workspaceResourceDesc = "List all workspaces"
	// 查询工作区文件清单
	workspaceFileResourceName = "Workspace_file_list"
	// 工作区文件清单，查看指定工作区文件清单
	workspaceFileResourceDesc = "Workspace file list, view the list of files in the specified workspace"
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

// extractIDFromFileResourceContentURI extracts the ID from the URI(workspaceFileResourceURI)
// 例如"http://mulife.vip/magicCodeAgent/workspace/test/file"，取出test
func (s *MCP) extractIDFromFileResourceContentURI(uri string) (ret string, err error) {
	parts := strings.Split(uri, "/")
	if len(parts) != 7 {
		log.Errorf("invalid url, url:%s", uri)
		err = errors.New("invalid url")
		return
	}
	if parts[6] != workspaceFileResourceName || parts[4] != workspaceResourceName {
		log.Errorf("invalid url, url:%s", uri)
		err = errors.New("invalid url")
		return
	}

	ret = parts[5]
	return
}

func (s *MCP) workspaceFileResourceContentURI(name string) string {
	urlVal, urlErr := url.Parse(mcpServerName)
	if urlErr != nil {
		log.Errorf("parse url failed, urlErr:%s", urlErr.Error())
		return ""
	}
	urlVal.Scheme = httpSchema
	urlVal.Path = path.Join(urlVal.Path, workspaceResourceName, name, workspaceFileResourceName)

	return urlVal.String()
}

func (s *MCP) resourceSet() {
	s.mcpServerPtr.AddResource(
		mcp.NewResource(
			s.workspaceResourceURI(),
			workspaceResourceName,
			mcp.WithMIMEType(mimeType),
			mcp.WithResourceDescription(workspaceResourceDesc),
		),
		s.workspaceResourceHandler,
	)

	s.mcpServerPtr.AddResourceTemplate(
		mcp.NewResourceTemplate(
			s.workspaceFileResourceURI(),
			workspaceFileResourceName,
			mcp.WithTemplateMIMEType(mimeType),
			mcp.WithTemplateDescription(workspaceFileResourceDesc),
		),
		s.workspaceFileResourceHandler,
	)
}

func (s *MCP) workspaceResourceHandler(
	ctx context.Context,
	request mcp.ReadResourceRequest) (ret []mcp.ResourceContents, err error) {
	ev := event.NewEvent(mc.FilterWorkspace, s.ID(), mc.WorkspaceModule, nil, nil)
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
	ev := event.NewEvent(mc.FilterWorkspaceFile, s.ID(), mc.WorkspaceModule, nil, idVal)
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
