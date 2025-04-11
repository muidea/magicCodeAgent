package common

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	fp "github.com/muidea/magicCommon/foundation/path"

	"github.com/muidea/magicToolkit/internal/config"
)

const MonitorModule = "/internal/modules/kernel/monitor"

const (
	AddWorkspace        = "/workspace/add"
	RemoveWorkspace     = "/workspace/remove"
	RefreshWorkspace    = "/workspace/refresh"
	EnableWorkspace     = "/workspace/enable"
	DisableWorkspace    = "/workspace/disable"
	FilterWorkspace     = "/workspace/filter"
	FilterWorkspaceFile = "/workspace/file/filter"
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
