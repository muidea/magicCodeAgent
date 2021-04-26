{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
{{range $mk, $mv := $childList}}
    {{$moduleName := $mv.Name}}
    {{range $ak, $av := $mv.Action}}
        func (s *impl) {{formatAction $serviceName $moduleName $av}}({{formatParam $serviceName $moduleName $av}})({{formatResult $serviceName $moduleName $av}}){
        }
    {{end}}
{{end}}

