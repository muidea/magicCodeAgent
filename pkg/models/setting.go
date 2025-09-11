package models

import (
	"path/filepath"

	fo "github.com/muidea/magicCommon/foundation/os"
)

type WorkspaceSetting struct {
	RootPath        string  `json:"rootPath"`
	VectThreshold   float32 `json:"vectThreshold"`
	HybridThreshold int     `json:"hybridThreshold"`
	TopK            int     `json:"topK"`
	ChunkSize       int     `json:"chunkSize"`
	Overlap         int     `json:"overlap"`
}

var defaultWorkspaceSetting = &WorkspaceSetting{
	RootPath:        filepath.Join(fo.GetHomeDir(), "codespace"),
	VectThreshold:   0.5,
	HybridThreshold: 1000,
	TopK:            5,
	ChunkSize:       512,
	Overlap:         30,
}

func DefaultWorkspaceSetting() *WorkspaceSetting {
	return defaultWorkspaceSetting
}
