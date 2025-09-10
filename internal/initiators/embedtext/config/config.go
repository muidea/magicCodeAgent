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

func GetBgeEmbedAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeEmbedAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetBgeAPIServer(), "bge/embed")
}

func GetBgeEmbedBulkAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeEmbedBulkAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetBgeAPIServer(), "bge/embed_bulk")
}

func GetBgeKeywordsAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("BgeKeywordsAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetBgeAPIServer(), "bge/keywords")
}
