package common

const InitatorModule = "/modules/initator"

const (
	CreateNewInitator = "/initator/create"
	RefreshInitator   = "/initator/refresh"
	DestroyInitator   = "/initator/destroy"
)

type CreateInitatorParam struct {
	InitatorName  string `json:"initatorName"`
	InitatorPath  string `json:"initatorPath"`
	ProjectName   string `json:"projectName"`
	WorkPath      string `json:"workPath"`
	EnableService bool   `json:"enableService"`
	EnableDao     bool   `json:"enableDao"`
}
