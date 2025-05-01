package config

import (
	"time"

	"github.com/muidea/magicCodeAgent/internal/config"
)

const deepseek = "deepseek"
const deepseekUrl = "https://api.deepseek.com"
const deepseekKey = ""
const maxRetries = 10
const retryInterval = 3 * time.Second

func GetDeepSeekUrl() string {
	models := config.Models()
	for _, model := range models {
		if model.Provider == deepseek {
			return deepseekUrl + model.BaseUrl
		}
	}
	return deepseekUrl
}

func GetDeepSeekAPIKey() string {
	models := config.Models()
	for _, model := range models {
		if model.Provider == deepseek {
			return model.APIKey
		}
	}
	return deepseekKey
}

func GetDeepSeekMaxRetries() int {
	return maxRetries
}

func GetDeepSeekRetryInterval() time.Duration {
	return retryInterval
}
