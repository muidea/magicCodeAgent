{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
{{range $mk, $mv := $childList}}
    {{$moduleName := $mv.Name}}
    {{$subModuleName := $mv.Name}}
    {{$jsonName := $mv.Lower $subModuleName}}
    type {{$subModuleName}}Result struct {
    commonDef.Result
    {{$subModuleName}} *model.{{$subModuleName}}View `json:"{{$jsonName}}"`
    }

    type {{$subModuleName}}ListResult struct {
    commonDef.Result
    Total    int64                 `json:"total"`
    {{$subModuleName}} []*model.{{$subModuleName}}View `json:"{{$jsonName}}"`
    }
{{end}}
