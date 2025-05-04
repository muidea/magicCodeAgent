package config

import (
	"time"

	"github.com/muidea/magicCodeAgent/internal/config"
)

// const deepseek = "deepseek"
const remoteDeepseekUrl = "https://api.deepseek.com"

// const localDeepseekUrl = "http://127.0.0.1:11434"
const deepseekKey = ""
const deepseekModelName = "deepseek-chat"
const maxRetries = 10
const retryInterval = 3 * time.Second

func GetDeepSeekUrl(name string) string {
	models := config.Models()
	for _, model := range models {
		if model.Name == name {
			return model.BaseUrl
		}
	}

	return remoteDeepseekUrl
}

func GetDeepSeekAPIKey(name string) string {
	models := config.Models()
	for _, model := range models {
		if model.Name == name {
			return model.APIKey
		}
	}
	return deepseekKey
}

func GetDeepSeekModelName(name string) string {
	models := config.Models()
	for _, model := range models {
		if model.Name == name {
			return model.Model
		}
	}
	return deepseekModelName
}

func GetDeepSeekMaxRetries() int {
	return maxRetries
}

func GetDeepSeekRetryInterval() time.Duration {
	return retryInterval
}
