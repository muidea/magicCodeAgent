package service

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
	"os"
	"os/exec"
	"path"
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

func (s *Service) declareServiceImport(moduleName string) (ret string) {
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
	"github.com/muidea/{{.ProjectName}}/core/module/{{toLower .ModuleName}}/biz"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleName": moduleName})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareObjectImport(moduleName string) (ret string) {
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
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/{{.ProjectName}}/common"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleName": moduleName})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `
type {{.ModuleName}} struct {
	routeRegistry     toolkit.RouteRegistry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry

	validator fu.Validator

	bizPtr *biz.{{.ModuleName}}
}

func New(
	{{toLower .ModuleName}}Biz *biz.{{.ModuleName}},
) *{{.ModuleName}} {
	ptr := &{{.ModuleName}}{
		bizPtr:    {{toLower .ModuleName}}Biz,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *{{.ModuleName}}) BindRegistry(
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
func (s *{{.ModuleName}}) RegisterRoute() { {{range .ItemsList}}
	s.roleRouteRegistry.AddHandler(common.Filter{{.ActionVal}}, "GET", cc.ReadPrivate, s.Filter{{.NameVal}})
	s.roleRouteRegistry.AddHandler(common.Query{{.ActionVal}}, "GET", cc.ReadPrivate, s.Query{{.NameVal}})
	s.roleRouteRegistry.AddHandler(common.Create{{.ActionVal}}, "POST", cc.WritePrivate, s.Create{{.NameVal}})
	s.roleRouteRegistry.AddHandler(common.Update{{.ActionVal}}, "PUT", cc.WritePrivate, s.Update{{.NameVal}})
	s.roleRouteRegistry.AddHandler(common.Delete{{.ActionVal}}, "DELETE", cc.DeletePrivate, s.Delete{{.NameVal}})
{{end}}
}

func (s *{{.ModuleName}}) getCurrentEntity(ctx context.Context) (ret *cc.EntityView, err error) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	authVal, ok := curSession.GetOption(session.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*cc.EntityView)
	return
}

func (s *{{.ModuleName}}) getCurrentNamespace(ctx context.Context, res http.ResponseWriter, req *http.Request) (ret string) {
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

func (s *{{.ModuleName}}) queryEntity(sessionInfo *session.SessionInfo, id int, namespace string) (ret *cc.EntityView) {
	ret = s.bizPtr.QueryEntity(sessionInfo, id, namespace)
	return
}

func (s *{{.ModuleName}}) writeLog(ctx context.Context, res http.ResponseWriter, req *http.Request, memo string) {
	address := fn.GetHTTPRemoteAddress(req)
	entityPtr, _ := s.getCurrentEntity(ctx)
	namespace := s.getCurrentNamespace(ctx, res, req)
	s.bizPtr.WriteLog(memo, address, entityPtr, namespace)
}`
	type Item struct {
		ActionVal string
		NameVal   string
	}

	moduleName := util.ToTitle(ptr.FileName)
	itemList := []*Item{}
	for _, val := range ptr.Objects {
		actionVal := val.Name
		if val.Name != moduleName {
			actionVal = fmt.Sprintf("%s%s", moduleName, actionVal)
		}
		nameVal := val.Name
		itemList = append(itemList, &Item{actionVal, nameVal})
	}

	factoryPtr := factory.NewMemoryFactor(defineDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "ItemsList": itemList})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Service) declareObjectFunc(moduleName string, ptr *parser.Object) (ret string) {
	funcDeclare := `{{$ModuleName := .ModuleName}}
{{range .ItemsList}}
func (s *{{$.ModuleName}}) filter{{.NameVal}}Lite(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	result := &common.{{.NameVal}}LiteListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}List, {{toLower .NameVal}}Total, {{toLower .NameVal}}Err := s.bizPtr.Filter{{.NameVal}}(filter, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		for _, val := range {{toLower .NameVal}}List {
			ptr := &common.{{.NameVal}}Lite{}
			ptr.From{{.NameVal}}(val)
			result.{{.NameVal}} = append(result.{{.NameVal}}, ptr)
		}

		result.Total = {{toLower .NameVal}}Total
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

func (s *{{$.ModuleName}}) filter{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request, filter *bc.QueryFilter) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.{{.NameVal}}ListResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}List, {{toLower .NameVal}}Total, {{toLower .NameVal}}Err := s.bizPtr.Filter{{.NameVal}}(filter, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		for _, val := range {{toLower .NameVal}}List {
			ptr := &common.{{.NameVal}}View{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			ptr.From{{.NameVal}}(val, entityPtr)
			result.{{.NameVal}} = append(result.{{.NameVal}}, ptr)
		}

		result.Total = {{toLower .NameVal}}Total
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

func (s *{{$.ModuleName}}) Filter{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	filter := fu.NewFilter()
	filter.Decode(req)
	queryFilter := bc.NewFilter()
	if filter.Pagination != nil {
		queryFilter.Page(filter.Pagination)
	} else {
		queryFilter.Page(fu.NewPagination(20, 1))
	}
	queryFilter.Sort(fu.NewSortFilter("UpdateTime", false))

	modelVal, modelOK := filter.Get("mode")
	if modelOK {
		if modelVal == "1" {
			s.filter{{.NameVal}}Lite(ctx, res, req, queryFilter)
			return
		}

		if modelVal == "2" {
			s.filter{{.NameVal}}(ctx, res, req, queryFilter)
			return
		}

		res.WriteHeader(http.StatusNotFound)
		return
	}

	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.{{.NameVal}}StatisticResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		summaryList := s.bizPtr.QuerySummary(common.{{.NameVal}}, curNamespace)
		for _, val := range summaryList {
			view := &common.TotalizerView{}
			view.FromTotalizer(val)
			result.Summary = append(result.Summary, view)
		}

		{{toLower .NameVal}}List, {{toLower .NameVal}}Total, {{toLower .NameVal}}Err := s.bizPtr.Filter{{.NameVal}}(queryFilter, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		for _, val := range {{toLower .NameVal}}List {
			view := &common.{{.NameVal}}View{}

			entityPtr := s.queryEntity(curSession.GetSessionInfo(), val.Creater, curNamespace)
			view.From{{.NameVal}}(val, entityPtr)
			result.{{.NameVal}} = append(result.{{.NameVal}}, view)
		}
		result.Total = {{toLower .NameVal}}Total
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

func (s *{{$.ModuleName}}) Query{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.{{.NameVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}Ptr, {{toLower .NameVal}}Err := s.bizPtr.Query{{.NameVal}}(id, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), {{toLower .NameVal}}Ptr.Creater, curNamespace)
		result.{{.NameVal}} = &common.{{.NameVal}}View{}
		result.{{.NameVal}}.From{{.NameVal}}({{toLower .NameVal}}Ptr, entityPtr)
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

func (s *{{$.ModuleName}}) Create{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{.NameVal}}Result{}
	for {
		param := &common.{{.NameVal}}Param{}
		err := fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}
		curEntity, curErr := s.getCurrentEntity(ctx)
		if curErr != nil {
			result.ErrorCode = cd.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}Ptr, {{toLower .NameVal}}Err := s.bizPtr.Create{{.NameVal}}(param, curEntity.ID, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		// TODO
		memo := fmt.Sprintf("新增{{.NameVal}}:%s", {{toLower .NameVal}}Ptr.Name)
		s.writeLog(ctx, res, req, memo)

		result.{{.NameVal}} = &common.{{.NameVal}}View{}
		result.{{.NameVal}}.From{{.NameVal}}({{toLower .NameVal}}Ptr, curEntity)
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

func (s *{{$.ModuleName}}) Update{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	result := &common.{{.NameVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		param := &common.{{.NameVal}}Param{}
		err = fn.ParseJSONBody(req, s.validator, param)
		if err != nil {
			result.ErrorCode = cd.IllegalParam
			result.Reason = "invalid param"
			break
		}

		curEntity, curErr := s.getCurrentEntity(ctx)
		if curErr != nil {
			result.ErrorCode = cd.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}Ptr, {{toLower .NameVal}}Err := s.bizPtr.Update{{.NameVal}}(id, param, curEntity.ID, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		memo := fmt.Sprintf("更新{{.NameVal}}:%s", {{toLower .NameVal}}Ptr.Name)
		s.writeLog(ctx, res, req, memo)

		result.{{.NameVal}} = &common.{{.NameVal}}View{}
		result.{{.NameVal}}.From{{.NameVal}}({{toLower .NameVal}}Ptr, curEntity)
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

func (s *{{$.ModuleName}}) Delete{{.NameVal}}(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	result := &common.{{.NameVal}}Result{}
	for {
		id, err := fn.SplitRESTID(req.URL.Path)
		if err != nil || id == 0 {
			result.ErrorCode = cd.Failed
			result.Reason = "invalid param"
			break
		}

		curNamespace := s.getCurrentNamespace(ctx, res, req)
		{{toLower .NameVal}}Ptr, {{toLower .NameVal}}Err := s.bizPtr.Delete{{.NameVal}}(id, curNamespace)
		if {{toLower .NameVal}}Err != nil {
			result.ErrorCode = cd.Failed
			result.Reason = {{toLower .NameVal}}Err.Error()
			break
		}

		memo := fmt.Sprintf("删除{{.NameVal}}:%s", {{toLower .NameVal}}Ptr.Name)
		s.writeLog(ctx, res, req, memo)

		entityPtr := s.queryEntity(curSession.GetSessionInfo(), {{toLower .NameVal}}Ptr.Creater, curNamespace)
		result.{{.NameVal}} = &common.{{.NameVal}}View{}
		result.{{.NameVal}}.FromArticle({{toLower .NameVal}}Ptr, entityPtr)
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
{{end}}`

	type Item struct {
		NameVal string
	}

	itemList := []*Item{}
	nameVal := ptr.Name
	itemList = append(itemList, &Item{nameVal})

	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemsList": itemList, "ModuleName": moduleName})
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

		_, err = fileHandle.WriteString(s.declareServiceImport(moduleName))
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
		err = s.createObject(moduleName, val)
		if err != nil {
			return
		}
	}

	return
}

func (s *Service) createObject(moduleName string, ptr *parser.Object) (err error) {
	commonFile := fmt.Sprintf("%s.go", util.ToLower(ptr.Name))
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(moduleName), "service", commonFile)

	log.Info(commonFilePath)
	func() {
		var fileHandle *os.File
		var fileErr error
		if moduleName != ptr.Name {
			fileHandle, fileErr = os.OpenFile(commonFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
			if fileErr != nil {
				err = fileErr
				return
			}
		} else {
			fileHandle, fileErr = os.OpenFile(commonFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if fileErr != nil {
				err = fileErr
				return
			}
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declarePackage())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectImport(moduleName))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectFunc(moduleName, ptr))
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
