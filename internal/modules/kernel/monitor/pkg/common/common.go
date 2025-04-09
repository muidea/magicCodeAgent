package common

import (
	cd "github.com/muidea/magicCommon/def"
)

const MonitorModule = "/internal/modules/kernel/monitor"

const (
	AddWorkSpace            = "/workspace/add"
	RemoveWorkSpace         = "/workspace/remove"
	RefreshWorkSpace        = "/workspace/refresh"
	FilterWorkSpaceFileList = "/workspace/file/filter"
)

type FileListResult struct {
	cd.Result
	Values []string `json:"values"`
}
