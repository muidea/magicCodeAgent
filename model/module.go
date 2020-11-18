package model

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
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

func (s *Module) UrlHandle(url string) (ret string) {
	var err error
	switch url {
	case "query/":
		ret, err = s.FilterHandle()
	case "query/:id":
		ret, err = s.QueryHandle()
	case "create/":
		ret, err = s.CreateHandle()
	case "update/:id":
		ret, err = s.UpdateHandle()
	case "delete/:id":
		ret, err = s.DeleteHandle()
	}

	if err != nil {
		log.Panicf("%s,err:%s", url, err.Error())
	}

	return
}

func (s *Module) FilterHandle() (ret string, err error) {
	templateVal := `
		{{$module := .Module}}
		{{$moduleName := $module.Name}}
		{{$lowName := $module.LowerName}}
		func (s *{{$moduleName}}) Filter{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
			filter := commonDef.NewFilter()
			filter.Decode(req)
		
			result := &common.{{$moduleName}}ListResult{ {{$moduleName}}: []*model.{{$moduleName}}View{} }
			for {
				{{$lowName}}List, {{$lowName}}Total, {{$lowName}}Err := s.{{$lowName}}Biz.Filter(filter)
				if {{$lowName}}Err != nil {
					result.ErrorCode = commonDef.Failed
					result.Reason = {{$lowName}}Err.Error()
					break
				}
		
				for _, val := range {{$lowName}}List {
					view := &model.{{$moduleName}}View{}
					view.From{{$moduleName}}(val)
		
					result.{{$moduleName}} = append(result.{{$moduleName}}, view)
				}
		
				result.Total = {{$lowName}}Total
				result.ErrorCode = commonDef.Success
				break
			}
		
			block, err := json.Marshal(result)
			if err == nil {
				res.Write(block)
				return
			}
		
			res.WriteHeader(http.StatusExpectationFailed)
		}
`
	templateHandle, templateErr := template.New("FilterHandle").Parse(templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	dataBuffer := &bytes.Buffer{}
	data := map[string]interface{}{"Module": s}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) QueryHandle() (ret string, err error) {
	templateVal := `
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.LowerName}}
func (s *{{$moduleName}}) Query{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("query {{$lowName}} failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.QueryByID(id)
		if {{$lowName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowName}}Ptr)
		result.{{$moduleName}} = view
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
`
	templateHandle, templateErr := template.New("QueryHandle").Parse(templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	dataBuffer := &bytes.Buffer{}
	data := map[string]interface{}{"Module": s}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) CreateHandle() (ret string, err error) {
	templateVal := `
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.LowerName}}
func (s *{{$moduleName}}) Create{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		{{$lowName}} := &model.{{$moduleName}}{}
		err := net.ParseJSONBody(req, {{$lowName}})
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curEntity,curErr := s.getCurrentEntity(res, req)
		if curErr != nil{
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = curErr.Error()
			break
		}

		curNamespace, curErr := s.getCurrentNamespace(res, req)
		if curErr != nil{
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = curErr.Error()
			break
		}

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Create({{$lowName}},curEntity,curNamespace)
		if {{$lowName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowName}}Ptr)
		result.{{$moduleName}} = view
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
`
	templateHandle, templateErr := template.New("CreateHandle").Parse(templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	dataBuffer := &bytes.Buffer{}
	data := map[string]interface{}{"Module": s}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) UpdateHandle() (ret string, err error) {
	templateVal := `
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.LowerName}}
func (s *{{$moduleName}}) Update{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Query()
		if {{$lowName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowName}}Ptr)
		result.{{$moduleName}} = view
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
`
	templateHandle, templateErr := template.New("UpdateHandle").Parse(templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	dataBuffer := &bytes.Buffer{}
	data := map[string]interface{}{"Module": s}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) DeleteHandle() (ret string, err error) {
	templateVal := `
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.LowerName}}
func (s *{{$moduleName}}) Delete{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("delete {{$lowName}} failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Delete(id)
		if {{$lowName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowName}}Ptr)

		result.{{$moduleName}} = view
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
`
	templateHandle, templateErr := template.New("DeleteHandle").Parse(templateVal)
	if templateErr != nil {
		err = templateErr
		return
	}

	dataBuffer := &bytes.Buffer{}
	data := map[string]interface{}{"Module": s}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}
