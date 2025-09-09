package models

type WorkspaceSetting struct {
	RootPath   string  `json:"rootPath"`
	Similarity float32 `json:"similarity"`
	TopK       int     `json:"topK"`
}
