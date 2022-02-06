package parser

type Define struct {
	ApiVersion  string    `json:"apiVersion"`
	ProjectName *string   `json:"projectName"`
	Module      []*Module `json:"module"`
	Model       *[]*Model `json:"model"`
}

type Module struct {
	ApiVersion string   `json:"apiVersion"`
	Name       string   `json:"name"`
	Desc       string   `json:"desc"`
	Icon       string   `json:"icon"`
	Model      []*Model `json:"model"`
}

type Model struct {
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Icon   string   `json:"icon"`
	Action []string `json:"action"`
}
