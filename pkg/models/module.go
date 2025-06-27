package models

type Define struct {
	ApiVersion  string    `json:"apiVersion"`
	ProjectName string    `json:"projectName"`
	ProjectType string    `json:"projectType"`
	Module      []*Module `json:"module"`
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
