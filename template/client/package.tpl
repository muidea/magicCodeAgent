{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$packageName := .PackageName}}
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

    BindSession(sessionInfo *commonCommon.SessionInfo)
    UnBindSession()

    Release()
}

// NewClient new client
func NewClient(serverURL string) Client {
    clnt := &impl{serverURL: serverURL, httpClient: &http.Client{}}

    return clnt
}

type impl struct {
    serverURL   string
    sessionInfo *commonCommon.SessionInfo
    contextInfo commonCommon.ContextInfo
    httpClient  *http.Client
}

func (s *impl) AttachContext(ctx commonCommon.ContextInfo) {
    s.contextInfo = ctx
}

func (s *impl) DetachContext() {
    s.contextInfo = nil
}

func (s *impl) getContextValues() url.Values {
    ret := url.Values{}
    if s.contextInfo != nil {
        ret = s.contextInfo.Encode(ret)
    }

    return ret
}

func (s *impl) BindSession(sessionInfo *commonCommon.SessionInfo) {
    s.sessionInfo = sessionInfo
}

func (s *impl) UnBindSession() {
    s.sessionInfo = nil
}

func (s *impl) Release() {
    if s.httpClient != nil {
        s.httpClient.CloseIdleConnections()
        s.httpClient = nil
    }
}
