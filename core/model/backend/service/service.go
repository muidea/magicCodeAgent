package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	log "github.com/cihub/seelog"

	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type Service struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Service {
	return &Service{rootPath: rootPath, projectName: projectName}
}

func (s *Service) declarePackage() (ret string) {
	packageDeclare := `
package service
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Service) declareServiceImport(ptr *parser.Package) (ret string) {
	importDeclare := `
import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	cc "github.com/muidea/magicCas/common"
	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/{{.ProjectName}}/common"
	"github.com/muidea/{{.ProjectName}}/core/module/{{formatModulePath .ModuleVal}}/biz"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareObjectImport(ptr *parser.Package, objectPtr *parser.Object) (ret string) {
	importDeclare := `
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	bc "github.com/muidea/magicBatis/common"

	cd "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/{{.ProjectName}}/common"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `{{$ModuleVal := .ModuleVal}}
type {{formatModuleName .ModuleVal}} struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

	bizPtr *biz.{{formatModuleName $ModuleVal}}
}

func New(
	{{formatModuleValue .ModuleVal}}Biz *biz.{{formatModuleName $ModuleVal}},
) *{{formatModuleName $ModuleVal}} {
	ptr := &{{formatModuleName $ModuleVal}}{
		bizPtr:    {{formatModuleValue .ModuleVal}}Biz,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *{{formatModuleName $ModuleVal}}) BindRegistry(
	routeRegistry toolkit.RouteRegistry,
	casRouteRegistry toolkit.CasRegistry,
	roleRouteRegistry toolkit.RoleRegistry) {

	s.routeRegistry = routeRegistry
	s.casRouteRegistry = casRouteRegistry
	s.roleRouteRegistry = roleRouteRegistry

	s.routeRegistry.SetApiVersion(common.ApiVersion)
	s.casRouteRegistry.SetApiVersion(common.ApiVersion)
	s.roleRouteRegistry.SetApiVersion(common.ApiVersion)
}

// RegisterRoute 注册路由
func (s *{{formatModuleName $ModuleVal}}) RegisterRoute() { {{with .ModuleVal}}{{range .Objects}}
	s.roleRouteRegistry.AddHandler(common.Filter{{formatObjectActionName $ModuleVal .}}, "GET", cc.ReadPrivate, s.Filter{{formatObjectName .}})
	s.roleRouteRegistry.AddHandler(common.Query{{formatObjectActionName $ModuleVal .}}, "GET", cc.ReadPrivate, s.Query{{formatObjectName .}})
	s.roleRouteRegistry.AddHandler(common.Create{{formatObjectActionName $ModuleVal .}}, "POST", cc.WritePrivate, s.Create{{formatObjectName .}})
	s.roleRouteRegistry.AddHandler(common.Update{{formatObjectActionName $ModuleVal .}}, "PUT", cc.WritePrivate, s.Update{{formatObjectName .}})
	s.roleRouteRegistry.AddHandler(common.Delete{{formatObjectActionName $ModuleVal .}}, "DELETE", cc.DeletePrivate, s.Delete{{formatObjectName .}})
{{end}}{{end}}
}

func (s *{{formatModuleName $ModuleVal}}) getCurrentEntity(ctx context.Context) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *{{formatModuleName $ModuleVal}}) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(cc.NamespaceID)
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

func (s *{{formatModuleName $ModuleVal}}) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *{{formatModuleName $ModuleVal}}) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	entityPtr, _ := s.getCurrentEntity(ctx)
	namespace := s.getCurrentNamespace(ctx, res, req)
	s.bizPtr.WriteLog(memo, address, entityPtr, namespace)
}`
	factoryPtr := factory.NewMemoryFactor(defineDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareObjectFunc(modulePtr *parser.Package, objectPtr *parser.Object) (ret string) {
	funcDeclare := `{{$ModuleVal := .ModuleVal}}{{$ObjectVal := .ObjectVal}}{{$HasEntity := .HasEntity}}
func (s *{{formatModuleName $ModuleVal}}) filter{{formatObjectName $ObjectVal}}Lite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.{{formatObjectName $ObjectVal}}LiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}List, {{formatObjectValue $ObjectVal}}Total, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Filter{{formatObjectName $ObjectVal}}(filter, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		for _, val := range {{formatObjectValue $ObjectVal}}List {
			ptr := &common.{{formatObjectName $ObjectVal}}Lite{}
			ptr.From{{formatObjectName $ObjectVal}}(val)
			result.{{formatObjectName $ObjectVal}} = append(result.{{formatObjectName $ObjectVal}}, ptr)
		}

		result.Total = {{formatObjectValue $ObjectVal}}Total
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) filter{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.{{formatObjectName $ObjectVal}}ListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}List, {{formatObjectValue $ObjectVal}}Total, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Filter{{formatObjectName $ObjectVal}}(filter, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		for _, val := range {{formatObjectValue $ObjectVal}}List {
			ptr := &common.{{formatObjectName $ObjectVal}}View{}
			ptr.From{{formatObjectName $ObjectVal}}(val)
			result.{{formatObjectName $ObjectVal}} = append(result.{{formatObjectName $ObjectVal}}, ptr)
		}

		result.Total = {{formatObjectValue $ObjectVal}}Total
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) Filter{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)
	queryFilter := bc.NewFilter()
	if filter.Pagination != nil {
		queryFilter.Page(filter.Pagination)
	} else {
		queryFilter.Page(fu.NewPagination(20, 1))
	}

	modelVal, modelOK := filter.Get("mode")
	if modelOK {
		if modelVal == common.LiteMode {
			s.filter{{formatObjectName $ObjectVal}}Lite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == common.ViewMode {
			s.filter{{formatObjectName $ObjectVal}}(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}
	result := &common.{{formatObjectName $ObjectVal}}StatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.Notify{{formatObjectActionName $ModuleVal $ObjectVal}}, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		{{formatObjectValue $ObjectVal}}List, {{formatObjectValue $ObjectVal}}Total, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Filter{{formatObjectName $ObjectVal}}(queryFilter, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		for _, val := range {{formatObjectValue $ObjectVal}}List {
			view := &common.{{formatObjectName $ObjectVal}}View{}
			view.From{{formatObjectName $ObjectVal}}(val)
			result.{{formatObjectName $ObjectVal}} = append(result.{{formatObjectName $ObjectVal}}, view)
		}
		result.Total = {{formatObjectValue $ObjectVal}}Total
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) Query{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{formatObjectName $ObjectVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}Ptr, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Query{{formatObjectName $ObjectVal}}(id, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		result.{{formatObjectName $ObjectVal}} = &common.{{formatObjectName $ObjectVal}}View{}
		result.{{formatObjectName $ObjectVal}}.From{{formatObjectName $ObjectVal}}({{formatObjectValue $ObjectVal}}Ptr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) Create{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{formatObjectName $ObjectVal}}Result{}
	for {
		param := &common.{{formatObjectName $ObjectVal}}Param{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}Ptr, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Create{{formatObjectName $ObjectVal}}(param, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		memo := fmt.Sprintf("新增{{formatObjectName $ObjectVal}}:%d", {{formatObjectValue $ObjectVal}}Ptr.ID)
		s.writeLog(ctx, res, req, memo)

		result.{{formatObjectName $ObjectVal}} = &common.{{formatObjectName $ObjectVal}}View{}
		result.{{formatObjectName $ObjectVal}}.From{{formatObjectName $ObjectVal}}({{formatObjectValue $ObjectVal}}Ptr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) Update{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{formatObjectName $ObjectVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.{{formatObjectName $ObjectVal}}Param{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}Ptr, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Update{{formatObjectName $ObjectVal}}(id, param, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		memo := fmt.Sprintf("更新{{formatObjectName $ObjectVal}}:%d", {{formatObjectValue $ObjectVal}}Ptr.ID)
		s.writeLog(ctx, res, req, memo)

		result.{{formatObjectName $ObjectVal}} = &common.{{formatObjectName $ObjectVal}}View{}
		result.{{formatObjectName $ObjectVal}}.From{{formatObjectName $ObjectVal}}({{formatObjectValue $ObjectVal}}Ptr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *{{formatModuleName $ModuleVal}}) Delete{{formatObjectName $ObjectVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{formatObjectName $ObjectVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{formatObjectValue $ObjectVal}}Ptr, {{formatObjectValue $ObjectVal}}Err := s.bizPtr.Delete{{formatObjectName $ObjectVal}}(id, curNamespace)
		if {{formatObjectValue $ObjectVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{formatObjectValue $ObjectVal}}Err.Error()
			break
		}

		memo := fmt.Sprintf("删除{{formatObjectName $ObjectVal}}:%d", {{formatObjectValue $ObjectVal}}Ptr.ID)
		s.writeLog(ctx, res, req, memo)

		result.{{formatObjectName $ObjectVal}} = &common.{{formatObjectName $ObjectVal}}View{}
		result.{{formatObjectName $ObjectVal}}.From{{formatObjectName $ObjectVal}}({{formatObjectValue $ObjectVal}}Ptr)
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
`

	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": modulePtr, "ObjectVal": objectPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := util.ToTitle(ptr.FileName)
	commonFile := fmt.Sprintf("%s.go", "service")
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(moduleName), "service", commonFile)

	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declarePackage())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareServiceImport(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareDefine(ptr))
		if err != nil {
			return
		}
	}()

	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", commonFilePath)
	err = cmd.Run()

	for _, val := range ptr.Objects {
		err = s.createObject(ptr, val)
		if err != nil {
			return
		}
	}

	return
}

func (s *Service) createObject(modulePtr *parser.Package, objectPtr *parser.Object) (err error) {
	commonFile := fmt.Sprintf("%s.go", util.ToLower(objectPtr.Name))
	commonFilePath := path.Join(s.rootPath, "core/module", mu.FormatModulePath(modulePtr), "service", commonFile)

	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declarePackage())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectImport(modulePtr, objectPtr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectFunc(modulePtr, objectPtr))
		if err != nil {
			return
		}
	}()

	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", commonFilePath)
	err = cmd.Run()
	return
}
