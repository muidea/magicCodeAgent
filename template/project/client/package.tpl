{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{range $modelList}}
        {{$modelName := .Name}}
        {{range $url := .Action}}
            {{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}){{end}}{{end}}
{{end}}

{{define "serviceFuncDefine"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{range $modelList}}
        {{$modelName := .Name}}
        {{range $url := .Action}}
            func (s *client) {{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}){
                return
            }{{end}}{{end}}
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
    AttachContext(ctx commonCommon.ContextInfo)
    DetachContext()

    {{range $serviceList}}
        {{template "serviceFuncDeclare" .}}{{end}}

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
    {{template "serviceFuncDefine" .}}{{end}}

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
