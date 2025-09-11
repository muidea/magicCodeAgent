package config

import (
	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

var workspaceSetting *models.WorkspaceSetting

func init() {
	workspaceSetting = config.WorkspaceSetting()
}

func GetVectThreshold() float32 {
	if workspaceSetting == nil || workspaceSetting.VectThreshold <= 0.0 {
		return models.DefaultWorkspaceSetting().VectThreshold
	}

	return workspaceSetting.VectThreshold
}

func GetHybridThreshold() int {
	if workspaceSetting == nil || workspaceSetting.HybridThreshold <= 0 {
		return models.DefaultWorkspaceSetting().HybridThreshold
	}
	return workspaceSetting.HybridThreshold
}

func GetTopK() int {
	if workspaceSetting == nil || workspaceSetting.TopK <= 0 {
		return models.DefaultWorkspaceSetting().TopK
	}
	return workspaceSetting.TopK
}

func GetChunkSize() int {
	if workspaceSetting == nil || workspaceSetting.ChunkSize <= 0 {
		return models.DefaultWorkspaceSetting().ChunkSize
	}
	return workspaceSetting.ChunkSize
}

func GetOverlap() int {
	if workspaceSetting == nil || workspaceSetting.Overlap <= 0 {
		return models.DefaultWorkspaceSetting().Overlap
	}
	return workspaceSetting.Overlap
}
