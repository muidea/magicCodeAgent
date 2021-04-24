{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
{{$serviceName}}
{{range $mk, $mv := $childList}}
    {{$subModuleName := $mv.Name}}
    {{$subModuleName}}
    {{$subModuleAction := $mv.Action}}
    {{range $ak, $av := $subModuleAction}}
        {{$av}}{{end}}{{end}}
