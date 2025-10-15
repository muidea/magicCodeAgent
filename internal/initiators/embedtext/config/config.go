package config

import (
	"fmt"
	"os"
)

var defaultNlpAPIServer = "https://api.mulife.vip"
var defaultNlpAPIUrl = "https://api.mulife.vip/nlp_service/api/v1/embedding/single"

func GetNLPAPIServer() string {
	apiServer, apiOK := os.LookupEnv("NlpAPIServer")
	if apiOK {
		return apiServer
	}

	return defaultNlpAPIServer
}

func GetNlpEmbedAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("NlpEmbedAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetNLPAPIServer(), "nlp_service/api/v1/embedding/single")
}

func GetNlpEmbedBatchAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("NlpEmbedBatchAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetNLPAPIServer(), "nlp_service/api/v1/embedding/batch")
}

func GetNlpKeywordsAPIUrl() string {
	apiVal, apiOK := os.LookupEnv("NlpKeywordsAPIUrl")
	if apiOK {
		return apiVal
	}

	return fmt.Sprintf("%s/%s", GetNLPAPIServer(), "nlp_service/api/v1/keywords/extract")
}
