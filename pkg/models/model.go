package models

// ModelInfo 定义AI模型配置结构
type ModelInfo struct {
	Name     string `json:"name"`     // 模型名称
	BaseUrl  string `json:"baseUrl"`  // 模型基础URL
	APIKey   string `json:"apiKey"`   // 模型API密钥
	Provider string `json:"provider"` // 模型提供商
	Model    string `json:"model"`
}
