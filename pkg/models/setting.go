package models

import (
	"path/filepath"

	fo "github.com/muidea/magicCommon/foundation/os"
)

type WorkspaceSetting struct {
	RootPath        string  `json:"rootPath"`
	VectSimilarity  float32 `json:"vectSimilarity"`
	HybridThreshold int     `json:"hybridThreshold"`
	TopK            int     `json:"topK"`
	ChunkSize       int     `json:"chunkSize"`
	Overlap         int     `json:"overlap"`
	MaxLoopCount    int     `json:"maxLoopCount"`
	EnableTrace     bool    `json:"enableTrace"`
}

var defaultWorkspaceSetting = &WorkspaceSetting{
	RootPath:        filepath.Join(fo.GetHomeDir(), "codespace"),
	VectSimilarity:  0.5,
	HybridThreshold: 1000,
	TopK:            10,
	ChunkSize:       512,
	Overlap:         30,
	MaxLoopCount:    100,
}

func DefaultWorkspaceSetting() *WorkspaceSetting {
	return defaultWorkspaceSetting
}
