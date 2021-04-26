{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "github.com/muidea/magicCommon/foundation/net"
    "github.com/muidea/magicCommon/foundation/util"

    commonCommon "github.com/muidea/magicCommon/common"
    commonDef "github.com/muidea/magicCommon/def"

    "github.com/muidea/{{$projectName}}/common"
    "github.com/muidea/{{$projectName}}/model"
)

// Client client interface
type Client interface {
    AttachContext(ctx commonCommon.ContextInfo)
    DetachContext()

    {{range $mk, $mv := $serviceList}}
        {{clientServiceFuncDeclare $mv}}{{end}}

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

{{range $mk, $mv := $serviceList}}
    {{clientServiceFuncDefine $mv}}{{end}}

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
