package models

import (
	"path/filepath"

	fo "github.com/muidea/magicCommon/foundation/os"
)

type WorkspaceSetting struct {
	RootPath        string  `json:"rootPath"`
	VectThreshold   float32 `json:"threshold"`
	HybridThreshold int     `json:"hybridThreshold"`
	TopK            int     `json:"topK"`
}

func DefaultWorkspaceSetting() *WorkspaceSetting {
	return &WorkspaceSetting{
		RootPath:        filepath.Join(fo.GetHomeDir(), "codespace"),
		VectThreshold:   0.5,
		HybridThreshold: 1000,
		TopK:            5,
	}
}
