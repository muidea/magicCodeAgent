package model

type Module struct {
	Name   string    `json:"name"`
	Desc   string    `json:"desc"`
	Action []string  `json:"action"`
	Child  []*Module `json:"child"`
}
