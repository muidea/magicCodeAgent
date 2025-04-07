package common

type CreateParam struct {
	Name          string `json:"name"`
	WorkingPath   string `json:"workingPath"`
	EnableService bool   `json:"enableService"`
	EnableDao     bool   `json:"enableDao"`
}
