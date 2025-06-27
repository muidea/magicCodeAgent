package common

const (
	LocalDeepseek  = "localDeepSeek"
	RemoteDeepseek = "remoteDeepSeek"
)

type DeepseekProvider struct{}

func (s *DeepseekProvider) String() string {
	return "DeepseekProvider"
}

const DefaultModel = "deepseek-chat"
const DefaultBaseUrl = "https://api.deepseek.com"
const DefaultProvider = "deepseek"

const Deepseek = "deepseek"
