package common

const ModuleModule = "/modules/module"

const (
	CreateNewModule = "/module/create"
	RefreshModule   = "/module/refresh"
	DestroyModule   = "/module/destroy"
)

type CreateModuleParam struct {
	ModuleName    string `json:"moduleName"`
	ModulePath    string `json:"modulePath"`
	ProjectName   string `json:"projectName"`
	WorkPath      string `json:"workPath"`
	EnableService bool   `json:"enableService"`
	EnableDao     bool   `json:"enableDao"`
}
