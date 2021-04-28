{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}

{{- define "filterFuncDeclare"}}
    {{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, err error)
{{- end}}
{{- define "filterFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, err error) {
        return
    }
{{- end}}

{{- define "queryFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "queryFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        return
    }
{{- end}}

{{- define "createFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "createFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
        return
    }
{{- end}}

{{- define "updateFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "updateFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
        return
    }
{{- end}}

{{- define "deleteFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "deleteFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        return
    }
{{- end}}

{{- define "otherFuncDeclare"}}
    {{.FuncName}}()
{{- end}}
{{- define "otherFuncDefine"}}
    func (s *{{toLower .ServiceName}}Handler){{.FuncName}}(){
        return
    }
{{- end}}

package {{$packageName}}

import (
"github.com/muidea/magicBatis/def"
"github.com/muidea/{{$projectName}}/model"
)

// {{$serviceName}}Handler {{toLower $serviceName}} handler
type {{$serviceName}}Handler interface {
    {{range $mv := $childList}}
        {{range $av := $mv.Action}}
            {{- $funcName := formatAction $serviceName $mv.Name $av -}}
            {{with mergeArgs $funcName $serviceName $mv}}
                {{- $action := extractAction $av -}}
                {{if eq $action "filter"}}
                    {{- template "filterFuncDeclare" . -}}
                {{else if eq $action "query"}}
                    {{- template "queryFuncDeclare" . -}}
                {{else if eq $action "create"}}
                    {{- template "createFuncDeclare" . -}}
                {{else if eq $action "update"}}
                    {{- template "updateFuncDeclare" . -}}
                {{else if eq $action "delete"}}
                    {{- template "deleteFuncDeclare" . -}}
                {{else}}
                    {{- template "otherFuncDeclare" . -}}
                {{end}}
            {{- end}}
        {{- end}}
    {{- end}}
}

type {{toLower $serviceName}}Handler struct {
    commonHandler
}

{{range $mv := $childList}}
    {{range $av := $mv.Action}}
        {{- $funcName := formatAction $serviceName $mv.Name $av -}}
        {{with mergeArgs $funcName $serviceName $mv}}
            {{- $action := extractAction $av -}}
            {{if eq $action "filter"}}
                {{- template "filterFuncDefine" . -}}
            {{else if eq $action "query"}}
                {{- template "queryFuncDefine" . -}}
            {{else if eq $action "create"}}
                {{- template "createFuncDefine" . -}}
            {{else if eq $action "update"}}
                {{- template "updateFuncDefine" . -}}
            {{else if eq $action "delete"}}
                {{- template "deleteFuncDefine" . -}}
            {{else}}
                {{- template "otherFuncDefine" . -}}
            {{end}}
        {{end}}
    {{- end}}
{{- end}}