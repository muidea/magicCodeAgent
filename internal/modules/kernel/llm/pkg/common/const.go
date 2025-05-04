package common

const (
	LocalDeepseek  = "localDeepSeek"
	RemoteDeepseek = "remoteDeepSeek"
)

type DeepseekProvider struct{}

func (s *DeepseekProvider) String() string {
	return "DeepseekProvider"
}

const Deepseek = "deepseek"
