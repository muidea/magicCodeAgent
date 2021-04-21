package model

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type Define struct {
	ApiVersion string    `json:"apiVersion"`
	Module     []*Module `json:"module"`
}

type Module struct {
	Name   string    `json:"name"`
	Desc   string    `json:"desc"`
	Action []string  `json:"action"`
	Child  []*Module `json:"child"`
}

func (s *Module) Lower(val string) string {
	firstLetter := val[0:1]
	otherStr := val[1:]
	return fmt.Sprintf("%s%s", strings.ToLower(firstLetter), otherStr)
}

func (s *Module) getUrl(name, url string) (ret string) {
	value := fmt.Sprintf("%s/%s", name, url)
	items := []string{}
	for _, val := range strings.Split(value, "/") {
		if val == "" || val == ":id" {
			continue
		}

		items = append(items, strings.Title(val))
	}

	return strings.Join(items, "")
}

func (s *Module) IsPackage() bool {
	return len(s.Action) == 0
}

func (s *Module) FormatUrl(apiVersion, url, parent string) (ret string) {
	path := s.Name
	if parent != "" && parent != path {
		path = fmt.Sprintf("%s/%s", parent, path)
	}

	switch url {
	case "query/":
		ret = fmt.Sprintf("%sFilterURL = \"/%s/%s/query/\"", strings.Title(s.Name), apiVersion, strings.ToLower(path))
	case "query/:id":
		ret = fmt.Sprintf("%sQueryURL = \"/%s/%s/query/:id\"", strings.Title(s.Name), apiVersion, strings.ToLower(path))
	case "create/":
		ret = fmt.Sprintf("%sCreateURL = \"/%s/%s/create/\"", strings.Title(s.Name), apiVersion, strings.ToLower(path))
	case "update/:id":
		ret = fmt.Sprintf("%sUpdateURL = \"/%s/%s/update/:id\"", strings.Title(s.Name), apiVersion, strings.ToLower(path))
	case "delete/:id":
		ret = fmt.Sprintf("%sDeleteURL = \"/%s/%s/delete/:id\"", strings.Title(s.Name), apiVersion, strings.ToLower(path))
	default:
		ret = fmt.Sprintf("%sURL = \"/%s/%s/%s\"", s.getUrl(s.Name, url), apiVersion, strings.ToLower(path), strings.ToLower(url))
	}

	return
}

func (s *Module) UrlHandle(serviceName, url string) (ret string) {
	var err error
	switch url {
	case "query/":
		ret, err = s.FilterUrlHandle(serviceName)
	case "query/:id":
		ret, err = s.QueryUrlHandle(serviceName)
	case "create/":
		ret, err = s.CreateUrlHandle(serviceName)
	case "update/:id":
		ret, err = s.UpdateUrlHandle(serviceName)
	case "delete/:id":
		ret, err = s.DeleteUrlHandle(serviceName)
	}

	if err != nil {
		log.Panicf("%s,err:%s", url, err.Error())
	}

	return
}

func (s *Module) FilterUrlHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowServiceName := $module.Lower $serviceName}}
{{$lowModuleName := $module.Lower $moduleName}}
func (s *{{$serviceName}}) Filter{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	filter := commonDef.NewFilter()
	filter.Decode(req)

	result := &common.{{$moduleName}}ListResult{ {{$moduleName}}: []*model.{{$moduleName}}View{} }
	for {
		{{$lowModuleName}}List, {{$lowModuleName}}Total, {{$lowModuleName}}Err := s.{{$lowServiceName}}Biz.Filter{{$moduleName}}(filter)
		if {{$lowModuleName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowModuleName}}Err.Error()
			break
		}

		for _, val := range {{$lowModuleName}}List {
			view := &model.{{$moduleName}}View{}
			view.From{{$moduleName}}(val)

			result.{{$moduleName}} = append(result.{{$moduleName}}, view)
		}

		result.Total = {{$lowModuleName}}Total
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
	templateHandle, templateErr := template.New("FilterUrlHandle").Parse(templateVal)
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

func (s *Module) QueryUrlHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowServiceName := $module.Lower $serviceName}}
{{$lowModuleName := $module.Lower $moduleName}}
func (s *{{$serviceName}}) Query{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("query {{$lowModuleName}} failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		{{$lowModuleName}}Ptr, {{$lowModuleName}}Err := s.{{$lowServiceName}}Biz.Query{{$moduleName}}ByID(id)
		if {{$lowModuleName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowModuleName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowModuleName}}Ptr)
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
	templateHandle, templateErr := template.New("QueryUrlHandle").Parse(templateVal)
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

func (s *Module) CreateUrlHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowServiceName := $module.Lower $serviceName}}
{{$lowModuleName := $module.Lower $moduleName}}
func (s *{{$serviceName}}) Create{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		{{$lowModuleName}}Ptr := &model.{{$moduleName}}{}
		err := net.ParseJSONBody(req, {{$lowModuleName}}Ptr)
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

		curNamespace := s.getCurrentNamespace(res, req)
		{{$lowModuleName}}Ptr, {{$lowModuleName}}Err := s.{{$lowServiceName}}Biz.Create{{$moduleName}}({{$lowModuleName}}Ptr,curEntity,curNamespace)
		if {{$lowModuleName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowModuleName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowModuleName}}Ptr)
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
	templateHandle, templateErr := template.New("CreateUrlHandle").Parse(templateVal)
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

func (s *Module) UpdateUrlHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}	
{{$moduleName := $module.Name}}
{{$lowServiceName := $module.Lower $serviceName}}
{{$lowModuleName := $module.Lower $moduleName}}
func (s *{{$serviceName}}) Update{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("update {{$lowModuleName}} failed, illegal id, id:%s", value)
			result.ErrorCode = commonDef.Failed
			result.Reason = "invalid param"
			break
		}

		{{$lowModuleName}}Ptr := &model.{{$moduleName}}{}
		err = net.ParseJSONBody(req, {{$lowModuleName}}Ptr)
		if err != nil {
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = "invalid param"
			break
		}
		{{$lowModuleName}}Ptr.ID = id

		curEntity,curErr := s.getCurrentEntity(res, req)
		if curErr != nil{
			result.ErrorCode = commonDef.IllegalParam
			result.Reason = curErr.Error()
			break
		}

		curNamespace := s.getCurrentNamespace(res, req)
		{{$lowModuleName}}Ptr, {{$lowModuleName}}Err := s.{{$lowServiceName}}Biz.Update{{$moduleName}}({{$lowModuleName}}Ptr,curEntity,curNamespace)
		if {{$lowModuleName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowModuleName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowModuleName}}Ptr)
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
	templateHandle, templateErr := template.New("UpdateUrlHandle").Parse(templateVal)
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

func (s *Module) DeleteUrlHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowServiceName := $module.Lower $serviceName}}
{{$lowModuleName := $module.Lower $moduleName}}
func (s *{{$serviceName}}) Delete{{$moduleName}}(res http.ResponseWriter, req *http.Request) {
	result := &common.{{$moduleName}}Result{}
	for {
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil || id == 0 {
			log.Printf("delete {{$lowModuleName}} failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		{{$lowModuleName}}Ptr, {{$lowModuleName}}Err := s.{{$lowServiceName}}Biz.Delete{{$moduleName}}(id)
		if {{$lowModuleName}}Err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = {{$lowModuleName}}Err.Error()
			break
		}

		view := &model.{{$moduleName}}View{}
		view.From{{$moduleName}}({{$lowModuleName}}Ptr)

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
	templateHandle, templateErr := template.New("DeleteUrlHandle").Parse(templateVal)
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

func (s *Module) RouteHandle(serviceName, url string) (ret string) {
	var err error
	switch url {
	case "query/":
		ret, err = s.FilterRouteHandle(serviceName)
	case "query/:id":
		ret, err = s.QueryRouteHandle(serviceName)
	case "create/":
		ret, err = s.CreateRouteHandle(serviceName)
	case "update/:id":
		ret, err = s.UpdateRouteHandle(serviceName)
	case "delete/:id":
		ret, err = s.DeleteRouteHandle(serviceName)
	}

	if err != nil {
		log.Panicf("%s,err:%s", url, err.Error())
	}

	return
}

func (s *Module) FilterRouteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $serviceName}}
s.privateRouteRegistry.AddHandler(common.Filter{{$moduleName}}URL, "GET", casModel.ReadPrivate, s.Filter{{$moduleName}})
`
	templateHandle, templateErr := template.New("FilterRouteHandle").Parse(templateVal)
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

func (s *Module) QueryRouteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $serviceName}}
s.privateRouteRegistry.AddHandler(common.Query{{$moduleName}}URL, "GET", casModel.ReadPrivate, s.Query{{$moduleName}})
`
	templateHandle, templateErr := template.New("QueryRouteHandle").Parse(templateVal)
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

func (s *Module) CreateRouteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $serviceName}}
s.privateRouteRegistry.AddHandler(common.Create{{$moduleName}}URL, "POST", casModel.WritePrivate, s.Create{{$moduleName}})
`
	templateHandle, templateErr := template.New("CreateRouteHandle").Parse(templateVal)
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

func (s *Module) UpdateRouteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $serviceName}}
s.privateRouteRegistry.AddHandler(common.Update{{$moduleName}}URL, "PUT", casModel.WritePrivate, s.Update{{$moduleName}})
`
	templateHandle, templateErr := template.New("UpdateRouteHandle").Parse(templateVal)
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

func (s *Module) DeleteRouteHandle(serviceName string) (ret string, err error) {
	templateVal := `
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $serviceName}}
s.privateRouteRegistry.AddHandler(common.Delete{{$moduleName}}URL, "DELETE", casModel.DeletePrivate, s.Delete{{$moduleName}})
`
	templateHandle, templateErr := template.New("DeleteRouteHandle").Parse(templateVal)
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
