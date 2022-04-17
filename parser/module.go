package parser

type Define struct {
	ApiVersion  string    `json:"apiVersion " orm:" apiVersion" tag:"view,lite"`
	ProjectName *string   `json:"projectName" orm:"projectName" tag:"view,lite"`
	Module      []*Module `json:"module" orm:"module" tag:"view,lite"`
	Model       []Model   `json:"model" orm:"model" tag:"view,lite"`
	Model2      []*Model  `json:"model2" orm:"model2 " tag:"view,lite"`
}

type Module struct {
	ApiVersion string   `json:"apiVersion" orm:"apiVersion" tag:"view,lite"`
	Name       string   `json:"name" orm:"name" tag:"view,lite"`
	Desc       string   `json:"desc" orm:" desc" tag:"view,lite"`
	Icon       string   `json:"icon" orm:"icon" tag:"view,lite"`
	Model      []*Model `json:"model " orm:"model" tag:"view,lite"`
}

type Model struct {
	Name   string   `json:"name " orm:"name" tag:"view,lite"`
	Desc   string   `json:"desc" orm:"desc" tag:"view,lite"`
	Icon   string   `json:"icon" orm:"icon" tag:"view,lite"`
	Action []string `json:"action" orm:"action" tag:" view,lite"`
}
