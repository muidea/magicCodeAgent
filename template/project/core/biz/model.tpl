{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := $model.Name}}
package {{$packageName}}

import (
"github.com/muidea/magicBatis/def"
"github.com/muidea/magicCommon/event"
"github.com/muidea/magicCommon/task"

"github.com/muidea/{{$projectName}}/model"

"github.com/muidea/{{$projectName}}/core/handler"
)

{{- define "filterFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, err error) {
    return
    }
{{- end}}

{{- define "queryFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}

{{- define "createFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}

{{- define "updateFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}

{{- define "deleteFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}

{{- define "otherFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(){
    return
    }
{{- end}}

{{if eq $serviceName $modelName}}
type {{$serviceName}} struct {
{{toLower $serviceName}}Handler handler.{{$serviceName}}Handler
}

func New(
{{toLower $serviceName}}Handler handler.{{$serviceName}}Handler,
eventHub event.Hub,
backgroundRoutine *task.BackgroundRoutine,
) *{{$serviceName}} {
    {{toLower $serviceName}} := &{{$serviceName}}{}

    return {{toLower $serviceName}}
}

{{end}}


{{range $av := $model.Action}}
    {{- $funcName := formatAction $serviceName $model.Name $av -}}
    {{with mergeArgs $funcName $serviceName $model}}
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
