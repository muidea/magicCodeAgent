{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := .Model.Name}}
package {{$packageName}}

import (
"encoding/json"
"fmt"
"net"
"net/http"
"strconv"
"strings"

log "github.com/cihub/seelog"

fn "github.com/muidea/magicCommon/foundation/net"
fu "github.com/muidea/magicCommon/foundation/util"

commonCommon "github.com/muidea/magicCommon/common"
commonDef "github.com/muidea/magicCommon/def"
commonSession "github.com/muidea/magicCommon/session"
engine "github.com/muidea/magicEngine"

casCommon "github.com/muidea/magicCas/common"
casModel "github.com/muidea/magicCas/model"
casToolkit "github.com/muidea/magicCas/toolkit"

"github.com/muidea/{{$projectName}}/common"
"github.com/muidea/{{$projectName}}/core/biz/{{toLower $serviceName}}"
"github.com/muidea/{{$projectName}}/core/dao"
"github.com/muidea/{{$projectName}}/model"
)


{{- define "filterFuncDefine"}}
    func (s *{{.ServiceName}}){{toCamelCase .FuncName}}Lite(res http.ResponseWriter, req *http.Request, filter *commonDef.Filter) {
        result := &common.{{.Model.Name}}ListResult{ {{.Model.Name}}: []*common.{{.Model.Name}}Lite{} }
        for {
            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}List, _, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(filter, curNamespace)
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
        filter := commonDef.NewFilter()
        filter.Decode(req)
    
        modelVal, modelOK := filter.Get("mode")
        if modelOK && modelVal == "1" {
            s.{{toCamelCase .FuncName}}Lite(res, req, filter)
            return
        }

        result := &common.{{.Model.Name}}StatisticResult{ {{.Model.Name}}: []*model.{{.Model.Name}}{} }
        for {
            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}List, {{toLower .Model.Name}}Total, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(filter, curNamespace)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            for _, val := range {{toLower .Model.Name}}List {
                view := &common.{{.Model.Name}}View{}
                view.FromPartner(val)
                result.{{.Model.Name}} = append(result.{{.Model.Name}},view)
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
            _, value := fn.SplitRESTAPI(req.URL.Path)
            id, err := strconv.Atoi(value)
            if err != nil || id == 0 {
                result.ErrorCode = commonDef.Failed
                result.Reason = "invalid param"
                break
            }

            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(id, curNamespace)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            result.{{.Model.Name}} = &common.{{.Model.Name}}View{}
            result.{{.Model.Name}}.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
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
    func (s *{{.ServiceName}}){{.FuncName}}(res http.ResponseWriter, req *http.Request) {
        result := &common.{{.Model.Name}}Result{}
        for {
            param := &common.{{.Model.Name}}Param{}
            err := fn.ParseJSONBody(req, s.validator, param)
            if err != nil {
                result.ErrorCode = commonDef.IllegalParam
                result.Reason = "invalid param"
                break
            }
            curEntity, curErr := s.getCurrentEntity(res, req)
            if curErr != nil {
                result.ErrorCode = commonDef.InvalidAuthority
                result.Reason = "invalid authority"
                break
            }
            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(param, curEntity.ID, curNamespace)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            // TODO
            memo := fmt.Sprintf("新建{{.Model.Name}}:%s", {{toLower .Model.Name}}Ptr)
            s.writeLog(res, req, memo)

            result.{{.Model.Name}} = &common.{{.Model.Name}}View{}
            result.{{.Model.Name}}.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
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
            _, value := fn.SplitRESTAPI(req.URL.Path)
            id, err := strconv.Atoi(value)
            if err != nil || id == 0 {
                log.Errorf("update partner failed, illegal id, id:%s", value)
                result.ErrorCode = commonDef.Failed
                result.Reason = "invalid param"
                break
            }

            param := &model.{{.Model.Name}}{}
            err = fn.ParseJSONBody(req, s.validator, param)
            if err != nil {
                result.ErrorCode = commonDef.IllegalParam
                result.Reason = "invalid param"
                break
            }
            param.ID = id
            curEntity, curErr := s.getCurrentEntity(res, req)
            if curErr != nil {
                result.ErrorCode = commonDef.InvalidAuthority
                result.Reason = "invalid authority"
                break
            }

            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(param, curEntity.ID, curNamespace)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            // TODO
            memo := fmt.Sprintf("更新{{.Model.Name}}:%s", {{toLower .Model.Name}}Ptr)
            s.writeLog(res, req, memo)

            result.{{.Model.Name}} = &common.{{.Model.Name}}View{}
            result.{{.Model.Name}}.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
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
            _, value := fn.SplitRESTAPI(req.URL.Path)
            id, err := strconv.Atoi(value)
            if err != nil || id == 0 {
                result.ErrorCode = commonDef.Failed
                result.Reason = "invalid param"
                break
            }

            curNamespace := s.getCurrentNamespace(res,req)
            {{toLower .Model.Name}}Ptr, {{toLower .Model.Name}}Err := s.{{toLower .ServiceName}}Biz.{{.FuncName}}(id, curNamespace)
            if {{toLower .Model.Name}}Err != nil {
                result.ErrorCode = commonDef.Failed
                result.Reason = {{toLower .Model.Name}}Err.Error()
                break
            }

            // TODO
            memo := fmt.Sprintf("删除{{.Model.Name}}:%s", {{toLower .Model.Name}}Ptr)
            s.writeLog(res, req, memo)

            result.{{.Model.Name}} = &common.{{.Model.Name}}View{}
            result.{{.Model.Name}}.From{{.Model.Name}}({{toLower .Model.Name}}Ptr)
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
        roleRouteRegistry casToolkit.RoleRegistry
        validator fu.Validator

        common dao.Common
        {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}}
    }

    func New(
    common dao.Common,
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}},
    ) *{{$serviceName}} {
        ptr := &{{$serviceName}}{
            common: common,
            {{toLower $serviceName}}Biz: {{toLower $serviceName}}Biz,
            validator: fu.NewFormValidator(),
        }

        return ptr
    }

    func (s *{{$serviceName}}) BindRegistry(sessionRegistry commonSession.Registry,roleRouteRegistry casToolkit.RoleRegistry) {
        s.sessionRegistry = sessionRegistry
        s.roleRouteRegistry = roleRouteRegistry
    }

    // Handle middleware handler
    func (s *{{$serviceName}}) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
        curSession := s.sessionRegistry.GetSession(res, req)

        sessionInfo := curSession.GetSessionInfo()
        sessionInfo.Scope = commonCommon.ShareSession

        values := req.URL.Query()
        values = sessionInfo.Encode(values)
        req.URL.RawQuery = values.Encode()

        ctx.Next()
    }

    // RegisterHandler 注册路由
    func (s *{{$serviceName}}) RegisterHandler(router engine.Router) {
        {{- range $av := $model.Action -}}
            {{$funcName := formatAction $serviceName $modelName $av}}
            s.roleRouteRegistry.AddHandler(common.{{$funcName}}URL, "{{getActionMethod $av}}", {{getActionPrivate $av}}, s.{{$funcName}})
        {{- end}}
    }


    func (s *{{$serviceName}}) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casModel.Entity, err error) {
        curSession := s.sessionRegistry.GetSession(res, req)
        authVal, ok := curSession.GetOption(commonCommon.AuthAccount)
        if !ok {
            err = fmt.Errorf("无效权限,未通过验证")
            return
        }
        ret = authVal.(*casModel.Entity)
        return
    }

    func (s *{{$serviceName}}) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
        namespace := req.Header.Get(casCommon.NamespaceID)
        if namespace != "" {
            ret = namespace
            return
        }

        items := strings.Split(req.Host, ":")
        if nil != net.ParseIP(items[0]) {
            ret = "default"
            return
        }

        items = strings.Split(items[0], ".")
        if len(items) <= 1 {
            ret = "default"
            return
        }

        ret = items[0]
        return
    }

    func (s *{{$serviceName}}) queryEntity(id int, sessionInfo *commonCommon.SessionInfo, namespace string) (ret *casModel.Entity) {
        entityPtr, entityErr := s.common.QueryEntity(id, sessionInfo, namespace)
        if entityErr != nil {
            ret = casModel.UnknownEntity()
            return
        }

        ret = entityPtr
        return
    }

    func (s *{{$serviceName}}) writeLog(res http.ResponseWriter, req *http.Request, memo string) {
        address := fn.GetHTTPRemoteAddress(req)
        curEntity, _ := s.getCurrentEntity(res, req)
        curNamespace := s.getCurrentNamespace(res, req)
        _, logErr := s.common.WriteLog(address,memo,curEntity.ID,curNamespace)
        if logErr != nil {
            log.Errorf("Write log failed, err:%s", logErr.Error())
        }
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
