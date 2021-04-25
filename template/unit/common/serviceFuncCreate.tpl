{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
{{range $mk, $mv := $childList}}
    {{$moduleName := $mv.Name}}
    func New{{$moduleName}}Param() *{{$moduleName}}Param {
        return &{{$moduleName}}Param{}
    }

    func New{{$moduleName}}Result() *{{$moduleName}}Result {
        return &{{$moduleName}}Result{}
    }

    func New{{$moduleName}}ListResult() *{{$moduleName}}ListResult {
        return &{{$moduleName}}ListResult{}
    }
{{end}}
