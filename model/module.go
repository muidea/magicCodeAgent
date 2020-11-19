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

func (s *Module) Lower(val string) string {
	return strings.ToLower(val)
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

func (s *Module) UrlHandle(serviceName, url string) (ret string) {
	var err error
	switch url {
	case "query/":
		ret, err = s.FilterHandle(serviceName)
	case "query/:id":
		ret, err = s.QueryHandle(serviceName)
	case "create/":
		ret, err = s.CreateHandle(serviceName)
	case "update/:id":
		ret, err = s.UpdateHandle(serviceName)
	case "delete/:id":
		ret, err = s.DeleteHandle(serviceName)
	}

	if err != nil {
		log.Panicf("%s,err:%s", url, err.Error())
	}

	return
}

func (s *Module) FilterHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
func (s *{{$serviceName}}) Filter{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	result := &common.{{$moduleName}}ListResult{ {{$moduleName}}: []*model.{{$moduleName}}View{} }
	for {
		{{$lowName}}List, {{$lowName}}Total, {{$lowName}}Err := s.{{$lowName}}Biz.Filter{{$moduleName}}(filter)
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
	data := map[string]interface{}{"Module": s, "ServiceName": serviceName}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) QueryHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
func (s *{{$serviceName}}) Query{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
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

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Query{{$moduleName}}ByID(id)
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
	data := map[string]interface{}{"Module": s, "ServiceName": serviceName}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) CreateHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
func (s *{{$serviceName}}) Create{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		{{$lowName}}Ptr := &model.{{$moduleName}}{}
		err := net.ParseJSONBody(req, {{$lowName}}Ptr)
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

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Create{{$moduleName}}({{$lowName}}Ptr,curEntity,curNamespace)
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
	data := map[string]interface{}{"Module": s, "ServiceName": serviceName}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) UpdateHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
func (s *{{$serviceName}}) Update{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("update {{$lowName}} failed, illegal id, id:%s", value)
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		{{$lowName}}Ptr := &model.{{$moduleName}}{}
		err = net.ParseJSONBody(req, {{$lowName}}Ptr)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}
		{{$lowName}}Ptr.ID = id

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

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Update{{$moduleName}}({{$lowName}}Ptr,curEntity,curNamespace)
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
	data := map[string]interface{}{"Module": s, "ServiceName": serviceName}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}

func (s *Module) DeleteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
func (s *{{$serviceName}}) Delete{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
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

		{{$lowName}}Ptr, {{$lowName}}Err := s.{{$lowName}}Biz.Delete{{$moduleName}}(id)
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
	data := map[string]interface{}{"Module": s, "ServiceName": serviceName}
	err = templateHandle.Execute(dataBuffer, data)
	if err != nil {
		return
	}
	ret = dataBuffer.String()

	return
}
