{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}

{{- define "filterFuncDeclare"}}
    {{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, total int64, err error)
{{- end}}
{{- define "filterFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, total int64, err error) {
        ret = []*model.{{.Model.Name}}{}
        total, err = s.batisClient.BatchQueryEntity(&ret, filter)
        return
    }
{{- end}}

{{- define "queryFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "queryFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        if id <= 0 {
            err = fmt.Errorf("illegal id value")
            return
        }
        ptr := &model.{{.Model.Name}}{ID:id}
        err = s.batisClient.QueryEntity(ptr)
        if err != nil {
            return
        }
        ret = ptr
        return
    }
{{- end}}

{{- define "createFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "createFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
        if ptr == nil {
            err = fmt.Errorf("illegal ptr value")
            return
        }
        err = s.batisClient.InsertEntity(ptr)
        if err != nil {
            return
        }
        ret = ptr
        return
    }
{{- end}}

{{- define "updateFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "updateFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
        err = s.batisClient.UpdateEntity(ptr)
        if err != nil {
            return
        }
        ret = ptr
        return
    }
{{- end}}

{{- define "deleteFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}
{{- define "deleteFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        if id <= 0 {
            err = fmt.Errorf("illegal id value")
            return
        }
        ptr := &model.{{.Model.Name}}{ID:id}
        err = s.batisClient.QueryEntity(ptr)
        if err != nil {
            return
        }
        err = s.batisClient.DeleteEntity(ptr)
        if err != nil {
            return
        }
        ret = ptr
        return
    }
{{- end}}

{{- define "otherFuncDeclare"}}
    {{.FuncName}}()
{{- end}}
{{- define "otherFuncDefine"}}
    func (s *{{toLower .ServiceName}}){{.FuncName}}(){
        return
    }
{{- end}}

package {{$packageName}}

import (
    "fmt"

    "github.com/muidea/magicBatis/def"
    "github.com/muidea/{{$projectName}}/model"
)

// {{$serviceName}} {{toLower $serviceName}} handler
type {{$serviceName}} interface {
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
                {{end}}
            {{- end}}
        {{- end}}
    {{- end}}
}

type {{toLower $serviceName}} struct {
    common
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
            {{end}}
        {{end}}
    {{- end}}
{{- end}}