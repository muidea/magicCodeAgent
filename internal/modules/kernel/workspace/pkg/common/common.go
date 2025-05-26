package common

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	fp "github.com/muidea/magicCommon/foundation/path"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

const WorkspaceModule = "/internal/modules/kernel/workspace"

const (
	QueryWorkspace                = "/workspace/query"
	CreateWorkspace               = "/workspace/create"
	DestroyWorkspace              = "/workspace/destroy"
	AddWorkspace                  = "/workspace/add"
	UpdateWorkspace               = "/workspace/update"
	RemoveWorkspace               = "/workspace/remove"
	RefreshWorkspace              = "/workspace/refresh"
	EnableWorkspace               = "/workspace/enable"
	DisableWorkspace              = "/workspace/disable"
	UpdateWorkspaceStage          = "/workspace/stage/update"
	FilterWorkspace               = "/workspace/filter"
	FilterWorkspaceFile           = "/workspace/file/filter"
	FilterWorkspaceFileTree       = "/workspace/file/tree/filter"
	QueryWorkspaceFileContent     = "/workspace/file/content/query"
	ChatWorkspaceContentSSE       = "/workspace/content/chat/sse"
	ChatWorkspaceContentMessage   = "/workspace/content/chat/message"
	ChatWorkspaceContentTerminate = "/workspace/content/chat/terminate"
	ChatWorkspaceContentHistory   = "/workspace/content/chat/history"
	FilterWorkspaceFlow           = "/workspace/flow/query"
	DebugWorkspaceContent         = "/workspace/content/debug"
)

type WorkspaceResult struct {
	cd.Result
	Value *config.Workspace `json:"value"`
}

type WorkspaceListResult struct {
	cd.Result
	Values []*config.Workspace `json:"values"`
}

type FileFilter struct {
	Workspace string
	Name      string
	Path      string
	Extension string
}

func (s *FileFilter) Parse(req *http.Request) {
	s.Workspace = req.URL.Query().Get("workspace")
	s.Name = req.URL.Query().Get("name")
	s.Path = req.URL.Query().Get("path")
	s.Extension = req.URL.Query().Get("extension")
}

func (s *FileFilter) Filter(val interface{}) bool {
	filePath, fileOK := val.(string)
	if !fileOK {
		return false
	}
	if !fp.Exist(filePath) {
		return false
	}

	// 文件名匹配检查
	if s.Name != "" {
		fileName := filepath.Base(filePath)
		pattern := strings.ReplaceAll(s.Name, "*", ".*")
		matched, _ := regexp.MatchString("(?i)"+pattern, fileName)
		if !matched {
			return false
		}
	}

	// 路径前缀匹配检查
	if s.Path != "" {
		if !strings.HasPrefix(filePath, s.Path) {
			return false
		}
	}

	// 扩展名匹配检查
	if s.Extension != "" {
		fileExt := strings.ToLower(filepath.Ext(filePath))
		if len(fileExt) > 0 && fileExt[0] == '.' {
			fileExt = fileExt[1:]
		}
		if fileExt != strings.ToLower(s.Extension) {
			return false
		}
	}

	return true
}

type FileListResult struct {
	cd.Result
	Values []string `json:"values"`
}

type DataNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*DataNode `json:"children,omitempty"`
}

type FileTreeListResult struct {
	cd.Result
	Values []*DataNode `json:"values"`
}

type FileContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type FileContentResult struct {
	cd.Result
	Value *FileContent `json:"value"`
}

type ChatMessage struct {
	Workspace string `json:"workspace"`
	Stage     string `json:"stage"`
	Content   string `json:"content"`
}

type ChatMessageResult struct {
	cd.Result
}

type ChatHistoryResult struct {
	cd.Result
	Values []*models.MessageInfo `json:"values"`
}

type FilterFlowResult struct {
	cd.Result
	Values []*models.Flow `json:"values"`
}

type DebugWorkspaceResult struct {
	cd.Result
}
