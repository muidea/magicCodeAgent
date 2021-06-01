package model

type Define struct {
	ApiVersion  string    `json:"apiVersion"`
	ProjectName string    `json:"projectName"`
	Module      []*Module `json:"module"`
}

type Module struct {
	ApiVersion string    `json:"apiVersion"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Action     []string  `json:"action"`
	Child      []*Module `json:"child"`
}
