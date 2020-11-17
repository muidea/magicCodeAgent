package model

import (
	"fmt"
	"strings"
)

type Module struct {
	Name   string    `json:"name"`
	Desc   string    `json:"desc"`
	Action []string  `json:"action"`
	Child  []*Module `json:"child"`
}

func (s *Module) LowerName() string {
	return strings.ToLower(s.Name)
}

func (s *Module) FormatUrl(url string) (ret string) {
	switch url {
	case "query/all/":
		ret = fmt.Sprintf("QueryAll%sURL = \"/api/v1/%s/query/all/\"", strings.Title(s.Name), strings.ToLower(s.Name))
	case "query/":
		ret = fmt.Sprintf("Filter%sURL = \"/api/v1/%s/query/\"", strings.Title(s.Name), strings.ToLower(s.Name))
	case "query/:id":
		ret = fmt.Sprintf("Query%sURL = \"/api/v1/%s/query/:id\"", strings.Title(s.Name), strings.ToLower(s.Name))
	case "create/":
		ret = fmt.Sprintf("Create%sURL = \"/api/v1/%s/create/\"", strings.Title(s.Name), strings.ToLower(s.Name))
	case "update/:id":
		ret = fmt.Sprintf("Update%sURL = \"/api/v1/%s/update/:id\"", strings.Title(s.Name), strings.ToLower(s.Name))
	case "delete/:id":
		ret = fmt.Sprintf("Delete%sURL = \"/api/v1/%s/delete/:id\"", strings.Title(s.Name), strings.ToLower(s.Name))
	}

	return
}
