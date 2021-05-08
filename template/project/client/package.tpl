{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

{{- define "filterFuncDefine"}}
    func (s *client){{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, total int64, err error) {
    return
    }
{{- end}}
{{- define "filterFuncDeclare"}}
    {{.FuncName}}(filter *def.QueryFilter) (ret []*model.{{.Model.Name}}, total int64, err error)
{{- end}}

{{- define "queryFuncDefine"}}
    func (s *client){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}
{{- define "queryFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "createFuncDefine"}}
    func (s *client){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}
{{- define "createFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "updateFuncDefine"}}
    func (s *client){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}
{{- define "updateFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "deleteFuncDefine"}}
    func (s *client){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
    return
    }
{{- end}}
{{- define "deleteFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "otherFuncDefine"}}
    func (s *client){{.FuncName}}(){
    return
    }
{{- end}}
{{- define "otherFuncDeclare"}}
    {{.FuncName}}()
{{- end}}

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{- range $mv := $modelList}}
        {{$modelName := $mv.Name}}
        {{- range $av := $mv.Action}}
            {{- $funcName := formatAction $serviceName $modelName $av -}}
            {{- with mergeArgs $funcName $serviceName $mv}}
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
{{end}}

{{define "serviceFuncDefine"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{- range $mv := $modelList}}
        {{$modelName := $mv.Name}}
        {{- range $av := $mv.Action}}
            {{- $funcName := formatAction $serviceName $modelName $av -}}
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
{{end}}


package {{$packageName}}

import (
"net/http"
"net/url"

"github.com/muidea/magicBatis/def"
commonCommon "github.com/muidea/magicCommon/common"

"github.com/muidea/{{$projectName}}/model"
)

// Client client interface
type Client interface {
    {{- range $serviceList}}
        {{template "serviceFuncDeclare" .}}
    {{- end}}

    AttachContext(ctx commonCommon.ContextInfo)
    DetachContext()

    BindSession(sessionInfo *commonCommon.SessionInfo)
    UnBindSession()

    Release()
}

// NewClient new client
func NewClient(serverURL string) Client {
    clnt := &client{serverURL: serverURL, httpClient: &http.Client{}}

    return clnt
}

type client struct {
    serverURL   string
    sessionInfo *commonCommon.SessionInfo
    contextInfo commonCommon.ContextInfo
    httpClient  *http.Client
}

{{range $serviceList}}
    {{template "serviceFuncDefine" .}}
{{- end}}

func (s *client) AttachContext(ctx commonCommon.ContextInfo) {
    s.contextInfo = ctx
}

func (s *client) DetachContext() {
    s.contextInfo = nil
}

func (s *client) getContextValues() url.Values {
    ret := url.Values{}
    if s.contextInfo != nil {
        ret = s.contextInfo.Encode(ret)
    }

    return ret
}

func (s *client) BindSession(sessionInfo *commonCommon.SessionInfo) {
    s.sessionInfo = sessionInfo
}

func (s *client) UnBindSession() {
    s.sessionInfo = nil
}

func (s *client) Release() {
    if s.httpClient != nil {
        s.httpClient.CloseIdleConnections()
        s.httpClient = nil
    }
}
