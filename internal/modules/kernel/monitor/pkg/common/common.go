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
	name      string
	path      string
	extension string
}

func (s *FileFilter) Parse(req *http.Request) {
	s.Workspace = req.URL.Query().Get("workspace")
	s.name = req.URL.Query().Get("name")
	s.path = req.URL.Query().Get("path")
	s.extension = req.URL.Query().Get("extension")
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
	if s.name != "" {
		fileName := filepath.Base(filePath)
		pattern := strings.ReplaceAll(s.name, "*", ".*")
		matched, _ := regexp.MatchString("(?i)"+pattern, fileName)
		if !matched {
			return false
		}
	}

	// 路径前缀匹配检查
	if s.path != "" {
		if !strings.HasPrefix(filePath, s.path) {
			return false
		}
	}

	// 扩展名匹配检查
	if s.extension != "" {
		fileExt := strings.ToLower(filepath.Ext(filePath))
		if len(fileExt) > 0 && fileExt[0] == '.' {
			fileExt = fileExt[1:]
		}
		if fileExt != strings.ToLower(s.extension) {
			return false
		}
	}

	return true
}

type FileListResult struct {
	cd.Result
	Values []string `json:"values"`
}
