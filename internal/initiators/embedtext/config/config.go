package config

import (
	"fmt"
	"os"
)

var defaultBgeAPIServer = "https://api.mulife.vip"
var defaultBgeAPIUrl = "https://api.mulife.vip/bge/embed"

func GetBgeAPIServer() string {
	apiServer, apiOK := os.LookupEnv("BgeAPIServer")
	if apiOK {
		return apiServer
	}

	return defaultBgeAPIServer
}

func GetBgeAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetBgeAPIServer(), "bge/embed")
}

func GetBgeAPIBulkUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeAPIBulkUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetBgeAPIServer(), "bge/embed_bulk")
}
