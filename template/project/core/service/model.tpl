{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := .Model.Name}}
package {{$packageName}}

import (
"encoding/json"
"net/http"
"strconv"

log "github.com/cihub/seelog"

"github.com/muidea/magicCommon/foundation/net"
"github.com/muidea/magicCommon/foundation/util"

commonDef "github.com/muidea/magicCommon/def"
commonSession "github.com/muidea/magicCommon/session"
engine "github.com/muidea/magicEngine"

"github.com/muidea/magicBatis/def"

"github.com/muidea/{{$projectName}}/common"
"github.com/muidea/{{$projectName}}/core/biz/{{toLower $serviceName}}"
"github.com/muidea/{{$projectName}}/core/dao"
"github.com/muidea/{{$projectName}}/model"
)


{{- define "filterFuncDefine"}}
    func (s *{{.ServiceName}})filter{{.Model.Name}}Lite(res http.ResponseWriter, req *http.Request, dataFilter *commonDef.Filter) {
        filter := def.NewFilter()
        filter.Page(dataFilter.PageFilter)

        result := &common.{{.Model.Name}}ListResult{ {{.Model.Name}}: []*common.{{.Model.Name}}Lite{} }
        for {
            {{toLower .Model.Name}}List, _, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Filter{{.Model.Name}}(filter)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            for _, val := range {{toLower .Model.Name}}List {
                {{toLower .Model.Name}}Lite := &common.{{.Model.Name}}Lite{}
                {{toLower .Model.Name}}Lite.From{{.Model.Name}}(val)
                result.{{.Model.Name}} = append(result.{{.Model.Name}}, {{toLower .Model.Name}}Lite)
            }

            result.ErrorCode = commonDef.Success
            break
        }

        block, err := json.Marshal(result)
        if err == nil {
            res.Write(block)
            return
        }

        res.WriteHeader(http.StatusExpectationFailed)
    }
    
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request) {
        dataFilter := commonDef.NewFilter()
        dataFilter.Decode(req)
    
        filter := def.NewFilter()
        filter.Page(dataFilter.PageFilter)
    
        modelVal, modelOK := dataFilter.Get("mode")
        if modelOK && modelVal == "1" {
            s.filter{{.Model.Name}}Lite(res, req, dataFilter)
            return
        }
    
        result := &common.{{.Model.Name}}StatisticResult{ {{.Model.Name}}: []*common.{{.Model.Name}}View{} }
        for {
        {{toLower .Model.Name}}List, {{toLower .Model.Name}}Total, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Filter{{.Model.Name}}(filter)
        if {{toLower .Model.Name}}Err != nil {
        result.ErrorCode = commonDef.Failed
        result.Reason = {{toLower .Model.Name}}Err.Error()
        break
        }
    
        for _, val := range {{toLower .Model.Name}}List {
        {{toLower .Model.Name}}View := &common.{{.Model.Name}}View{}
        {{toLower .Model.Name}}View.From{{.Model.Name}}(val)
        result.{{.Model.Name}} = append(result.{{.Model.Name}}, {{toLower .Model.Name}}View)
        }
    
        result.Total = {{toLower .Model.Name}}Total
        result.ErrorCode = commonDef.Success
        break
        }
    
        block, err := json.Marshal(result)
        if err == nil {
        res.Write(block)
        return
        }
    
        res.WriteHeader(http.StatusExpectationFailed)
    }
{{- end}}

{{- define "queryFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request) {
        result := &common.{{.Model.Name}}Result{}
        for {
            _, value := net.SplitRESTAPI(req.URL.Path)
            id, err := strconv.Atoi(value)
            if err != nil || id == 0 {
                result.ErrorCode = commonDef.Failed
                result.Reason = "invalid param"
                break
            }

            {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Query{{.Model.Name}}(id)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            {{toLower .Model.Name}}View := &common.{{.Model.Name}}View{}
            {{toLower .Model.Name}}View.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
            result.{{.Model.Name}} = {{toLower .Model.Name}}View
            result.ErrorCode = commonDef.Success
            break
        }

        block, err := json.Marshal(result)
        if err == nil {
        res.Write(block)
        return
        }

        res.WriteHeader(http.StatusExpectationFailed)
    }
{{- end}}

{{- define "createFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
        result := &common.{{.Model.Name}}Result{}
        for {
        param := &common.{{.Model.Name}}Param{}
        err := net.ParseJSONBody(req, s.validator, param)
        if err != nil {
        result.ErrorCode = commonDef.IllegalParam
        result.Reason = "invalid param"
        break
        }

        {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Create{{.Model.Name}}(param.To{{.Model.Name}}())
        if {{toLower .Model.Name}}Err != nil {
        result.ErrorCode = commonDef.Failed
        result.Reason = {{toLower .Model.Name}}Err.Error()
        break
        }

        {{toLower .Model.Name}}View := &common.{{.Model.Name}}View{}
        {{toLower .Model.Name}}View.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
        result.{{.Model.Name}} = {{toLower .Model.Name}}View
        result.ErrorCode = commonDef.Success
        break
        }

        block, err := json.Marshal(result)
        if err == nil {
        res.Write(block)
        return
        }

        res.WriteHeader(http.StatusExpectationFailed)
    }
{{- end}}

{{- define "updateFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
        result := &common.{{.Model.Name}}Result{}
        for {
        _, value := net.SplitRESTAPI(req.URL.Path)
        id, err := strconv.Atoi(value)
        if err != nil || id == 0 {
        log.Errorf("update partner failed, illegal id, id:%s", value)
        result.ErrorCode = commonDef.Failed
        result.Reason = "invalid param"
        break
        }

        param := &model.{{.Model.Name}}{}
        err = net.ParseJSONBody(req, s.validator, param)
        if err != nil {
        result.ErrorCode = commonDef.IllegalParam
        result.Reason = "invalid param"
        break
        }
        param.ID = id

        {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Update{{.Model.Name}}(param)
        if {{toLower .Model.Name}}Err != nil {
        result.ErrorCode = commonDef.Failed
        result.Reason = {{toLower .Model.Name}}Err.Error()
        break
        }

        {{toLower .Model.Name}}View := &common.{{.Model.Name}}View{}
        {{toLower .Model.Name}}View.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
        result.{{.Model.Name}} = {{toLower .Model.Name}}View
        result.ErrorCode = commonDef.Success
        break
        }

        block, err := json.Marshal(result)
        if err == nil {
        res.Write(block)
        return
        }

        res.WriteHeader(http.StatusExpectationFailed)
    }
{{- end}}

{{- define "deleteFuncDefine"}}
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request){
    result := &common.{{.Model.Name}}Result{}
    for {
    _, value := net.SplitRESTAPI(req.URL.Path)
    id, err := strconv.Atoi(value)
    if err != nil || id == 0 {
    result.ErrorCode = commonDef.Failed
    result.Reason = "invalid param"
    break
    }

    {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .Model.Name}}Biz.Delete{{.Model.Name}}(id)
    if {{toLower .Model.Name}}Err != nil {
    result.ErrorCode = commonDef.Failed
    result.Reason = {{toLower .Model.Name}}Err.Error()
    break
    }

    {{toLower .Model.Name}}View := &common.{{.Model.Name}}View{}
    {{toLower .Model.Name}}View.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
    result.{{.Model.Name}} = {{toLower .Model.Name}}View
    result.ErrorCode = commonDef.Success
    break
    }

    block, err := json.Marshal(result)
    if err == nil {
    res.Write(block)
    return
    }

    res.WriteHeader(http.StatusExpectationFailed)
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
