package config

import "os"

var defaultBgeAPIUrl = "https://api.mulife.vip/bge/embed"

func GetBgeAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeAPIUrl")
	if apiOK {
		return apiVal
	}

	return defaultBgeAPIUrl
}
