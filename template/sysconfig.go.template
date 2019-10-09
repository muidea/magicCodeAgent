package model

// SysConfig sysConfig
type SysConfig struct {
	ID    int          `json:"id" orm:"id key auto"`
	Name  string       `json:"name" orm:"name"`
	Items []ConfigItem `json:"item" orm:"item"`
}

// ConfigItem config item
type ConfigItem struct {
	//ID 唯一标示单元
	ID int `json:"id" orm:"id key auto"`
	// Name 名称
	Name string `json:"name" orm:"name"`
	//EMail 用户邮箱
	Value string `json:"value" orm:"value"`
}

// IsSame is same property
func (s *ConfigItem) IsSame(right *ConfigItem) bool {
	return s.Name == right.Name && s.Value == right.Value
}
