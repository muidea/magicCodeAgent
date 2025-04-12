package config

import (
	"path"

	"github.com/muidea/magicCodeAgent/internal/config"
)

const deepseek = "deepseek"
const deepseekUrl = "https://api.deepseek.com"
const deepseekKey = ""

func GetDeepSeekUrl() string {
	models := config.Models()
	for _, model := range models {
		if model.OwnedBy == deepseek {
			return path.Join(deepseekUrl, model.BaseUrl)
		}
	}
	return deepseekUrl
}

func GetDeepSeekAPIKey() string {
	models := config.Models()
	for _, model := range models {
		if model.OwnedBy == deepseek {
			return model.APIKey
		}
	}
	return deepseekKey
}
