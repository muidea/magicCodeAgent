{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}

{{define "filterFuncDeclare"}}
    {{.Name}}(filter *def.QueryFilter) (ret []*model.{{.Value.Name}}, err error)
{{end}}

{{define "queryFuncDeclare"}}
    {{.Name}}(id int) (ret *model.{{.Value.Name}}, err error)
{{end}}

{{define "createFuncDeclare"}}
    {{.Name}}(ptr *model.{{.Value.Name}}) (ret *model.{{.Value.Name}}, err error)
{{end}}

{{define "updateFuncDeclare"}}
    {{.Name}}(ptr *model.{{.Value.Name}}) (ret *model.{{.Value.Name}}, err error)
{{end}}

{{define "deleteFuncDeclare"}}
    {{.Name}}(id int) (ret *model.{{.Value.Name}}, err error)
{{end}}

{{define "otherFuncDeclare"}}
    {{.Name}}()
{{end}}

package {{$packageName}}

import (
"github.com/muidea/magicBatis/def"
"github.com/muidea/{{$projectName}}/model"
)

// {{$serviceName}}Handler {{toLower $serviceName}} handler
type {{$serviceName}}Handler interface {
    {{range $mv := $childList}}
        {{range $av := $mv.Action}}
            {{ $funcName := formatAction $serviceName $mv.Name $av}}
            {{with mergeArgs $funcName $mv }}
            {{$action := extractAction $av}}
            {{if eq $action "filter"}}
                {{template "filterFuncDeclare" .}}
            {{else if eq $action "query"}}
                {{template "queryFuncDeclare" .}}
            {{else if eq $action "create"}}
                {{template "createFuncDeclare" .}}
            {{else if eq $action "update"}}
                {{template "updateFuncDeclare" .}}
            {{else if eq $action "delete"}}
                {{template "deleteFuncDeclare" .}}
            {{else}}
                {{template "otherFuncDeclare" .}}
            {{end}}
            {{end}}
        {{end}}
    {{end}}
}

type {{toLower $serviceName}}Handler struct {
    commonHandler
}

{{range $mv := $childList}}
    {{$modelName := $mv.Name}}
    {{range $url := $mv.Action}}
        func (s *{{toLower $serviceName}}Handler){{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}){
            return
        }
    {{end}}
{{end}}