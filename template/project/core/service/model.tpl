{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := .Model.Name}}
package {{$packageName}}

import (
"net/http"

"github.com/muidea/magicCommon/foundation/util"
commonSession "github.com/muidea/magicCommon/session"
engine "github.com/muidea/magicEngine"

"github.com/muidea/{{$projectName}}/common"
"github.com/muidea/{{$projectName}}/core/biz/{{toLower $serviceName}}"
"github.com/muidea/{{$projectName}}/core/dao"
)


{{- define "filterFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request) {
    }
{{- end}}

{{- define "queryFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request) {
    }
{{- end}}

{{- define "createFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
    }
{{- end}}

{{- define "updateFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
    }
{{- end}}

{{- define "deleteFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
    }
{{- end}}

{{- define "otherFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
    }
{{- end}}

{{if eq $serviceName $modelName}}
    type {{$serviceName}} struct {
    sessionRegistry      commonSession.Registry
    validator util.Validator

    common dao.Common
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}}
    }

    func New(
    common dao.Common,
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}},
    ) *{{$serviceName}} {
    ptr := &{{$serviceName}}{
    common:common,
    {{toLower $serviceName}}Biz:{{toLower $serviceName}}Biz,
    validator:     util.NewFormValidator(),
    }

    return ptr
    }

    func (s *{{$serviceName}}) BindRegistry(sessionRegistry commonSession.Registry) {
        s.sessionRegistry = sessionRegistry
    }


    // Handle middleware handler
    func (s *{{$serviceName}}) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
        ctx.Next()
    }

    // RegisterHandler 注册路由
    func (s *{{$serviceName}}) RegisterHandler(router engine.Router) {
        {{range $av := $model.Action}}
            {{$funcName := formatAction $serviceName $modelName $av}}
            {{toVariable $funcName}}Route := engine.CreateRoute(common.{{$funcName}}URL, "{{getActionMethod $av}}", s.{{$funcName}})
            router.AddRoute({{toVariable $funcName}}Route, s)
        {{end}}
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
