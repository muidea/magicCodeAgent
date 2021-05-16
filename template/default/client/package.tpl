{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

{{- define "filterFuncDefine"}}
    func (s *client){{.FuncName}}(filter *commonDef.Filter) (ret []*model.{{.Model.Name}}, total int64, err error) {
        result := &common.{{.Model.Name}}StatisticResult{}

        vals := url.Values{}
        if s.sessionInfo != nil {
            vals = s.sessionInfo.Encode(vals)
        }
        if filter != nil {
            vals = filter.Encode(vals)
        }
        url, _ := url.ParseRequestURI(s.serverURL)
        url.Path = strings.Join([]string{url.Path, common.{{.FuncName}}URL}, "")
        url.RawQuery = vals.Encode()

        _, err = fn.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
        if err != nil {
            return
        }

        if result.ErrorCode != commonDef.Success {
        err = fmt.Errorf("查询{{.Model.Name}}失败，%s", result.Reason)
            return
        }

        ret = result.{{.Model.Name}}
        total = result.Total
        return
    }
{{- end}}
{{- define "filterFuncDeclare"}}
    {{.FuncName}}(filter *commonDef.Filter) (ret []*model.{{.Model.Name}}, total int64, err error)
{{- end}}

{{- define "queryFuncDefine"}}
    func (s *client){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        result := &common.{{.Model.Name}}Result{}

        vals := url.Values{}
        if s.sessionInfo != nil {
            vals = s.sessionInfo.Encode(vals)
        }
        url, _ := url.ParseRequestURI(s.serverURL)
        url.Path = strings.Join([]string{url.Path, common.{{.FuncName}}URL}, "")
        url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
        url.RawQuery = vals.Encode()

        _, err = fn.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
        if err != nil {
            return
        }

        if result.ErrorCode != commonDef.Success {
        err = fmt.Errorf("查询指定{{.Model.Name}}失败，%s", result.Reason)
            return
        }

        ret = result.{{.Model.Name}}
        return
    }
{{- end}}
{{- define "queryFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "createFuncDefine"}}
    func (s *client){{.FuncName}}(ptr *common.{{.Model.Name}}Param) (ret *model.{{.Model.Name}}, err error){
        result := &common.{{.Model.Name}}Result{}

        vals := url.Values{}
        if s.sessionInfo != nil {
            vals = s.sessionInfo.Encode(vals)
        }
        url, _ := url.ParseRequestURI(s.serverURL)
        url.Path = strings.Join([]string{url.Path, common.{{.FuncName}}URL}, "")
        url.RawQuery = vals.Encode()
        _, err = fn.HTTPPost(s.httpClient, url.String(), ptr, result, s.getContextValues())
        if err != nil {
            return
        }

        if result.ErrorCode != commonDef.Success {
        err = fmt.Errorf("新建{{.Model.Name}}失败，%s", result.Reason)
            return
        }

        ret = result.{{.Model.Name}}
        return
    }
{{- end}}
{{- define "createFuncDeclare"}}
    {{.FuncName}}(ptr *common.{{.Model.Name}}Param) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "updateFuncDefine"}}
    func (s *client){{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error){
        result := &common.{{.Model.Name}}Result{}

        vals := url.Values{}
        if s.sessionInfo != nil {
            vals = s.sessionInfo.Encode(vals)
        }
        url, _ := url.ParseRequestURI(s.serverURL)
        url.Path = strings.Join([]string{url.Path, common.{{.FuncName}}URL}, "")
        url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", ptr.ID))
        url.RawQuery = vals.Encode()

        _, err = fn.HTTPPut(s.httpClient, url.String(), ptr, result, s.getContextValues())
        if err != nil {
            return
        }

        if result.ErrorCode != commonDef.Success {
        err = fmt.Errorf("修改{{.Model.Name}}失败，%s", result.Reason)
            return
        }

        ret = result.{{.Model.Name}}
        return
    }
{{- end}}
{{- define "updateFuncDeclare"}}
    {{.FuncName}}(ptr *model.{{.Model.Name}}) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "deleteFuncDefine"}}
    func (s *client){{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error){
        result := &common.{{.Model.Name}}Result{}

        vals := url.Values{}
        if s.sessionInfo != nil {
            vals = s.sessionInfo.Encode(vals)
        }
        url, _ := url.ParseRequestURI(s.serverURL)
        url.Path = strings.Join([]string{url.Path, common.{{.FuncName}}URL}, "")
        url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
        url.RawQuery = vals.Encode()

        _, err = fn.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
        if err != nil {
            return
        }

        if result.ErrorCode != commonDef.Success {
        err = fmt.Errorf("删除{{.Model.Name}}失败，%s", result.Reason)
            return
        }

        ret = result.{{.Model.Name}}
        return
    }
{{- end}}
{{- define "deleteFuncDeclare"}}
    {{.FuncName}}(id int) (ret *model.{{.Model.Name}}, err error)
{{- end}}

{{- define "otherFuncDefine"}}
    func (s *client){{.FuncName}}(){
    // TODO
    return
    }
{{- end}}
{{- define "otherFuncDeclare"}}
    {{.FuncName}}()
{{- end}}

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{- range $mv := $modelList -}}
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
    {{- end -}}
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
"fmt"
"strings"
"net/http"
"net/url"

fn "github.com/muidea/magicCommon/foundation/net"

commonCommon "github.com/muidea/magicCommon/common"
commonDef "github.com/muidea/magicCommon/def"

"github.com/muidea/{{$projectName}}/common"
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
