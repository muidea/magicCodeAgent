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
	return workspaceSetting.VectThreshold
}

func GetHybridThreshold() int {
	return workspaceSetting.HybridThreshold
}

func GetTopK() int {
	return workspaceSetting.TopK
}
