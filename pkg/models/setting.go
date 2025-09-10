package models

type WorkspaceSetting struct {
	RootPath  string  `json:"rootPath"`
	Threshold float32 `json:"threshold"`
	TopK      int     `json:"topK"`
}
