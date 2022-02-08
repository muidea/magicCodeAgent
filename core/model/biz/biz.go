package biz

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

type Biz struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Biz {
	return &Biz{rootPath: rootPath, projectName: projectName}
}

func (s *Biz) declarePackage() (ret string) {
	packageDeclare := `
package biz
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Biz) declareBizImport(moduleName string) (ret string) {
	importDeclare := `
import (
	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/{{.ProjectName}}/common"
	"github.com/muidea/{{.ProjectName}}/core/base/biz"
	"github.com/muidea/{{.ProjectName}}/core/module/{{toLower .ModuleName}}/dao"
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

func (s *Biz) declareObjectImport(moduleName string) (ret string) {
	importDeclare := `
import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/{{.ProjectName}}/common"
	"github.com/muidea/{{.ProjectName}}/model"
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

func (s *Biz) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `
type {{.ModuleName}} struct {
	biz.Base

	{{toCamelCase .ModuleName}}Dao dao.{{.ModuleName}}

	endpointName string
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *Content {

	ptr := &{{.ModuleName}}{
		Base:         biz.New(common.{{.ModuleName}}Module, eventHub, backgroundRoutine),
		{{toCamelCase .ModuleName}}Dao:   dao.New(batisClient),
		endpointName: endpointName,
	}

	return ptr
}

func (s *{{.ModuleName}}) Notify(event event.Event, result event.Result) {
	return
}`
	moduleName := util.ToTitle(ptr.FileName)
	factoryPtr := factory.NewMemoryFactor(defineDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) declareObjectFunc(moduleName string, ptr *parser.Object) (ret string) {
	funcDeclare := `{{$ModuleName := .ModuleName}}
{{range .Items}}
func (s *{{$ModuleName}}) Filter{{formatNameVal $ModuleName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatNameVal $ModuleName .}}, total int64, err error) {
	ret, total, err = s.{{toCamelCase $ModuleName}}Dao.Filter{{formatNameVal $ModuleName .}}(filter, namespace)
	return
}

func (s *{{$ModuleName}}) Query{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	ret, err = s.{{toCamelCase $ModuleName}}Dao.Query{{formatNameVal $ModuleName .}}(id, namespace)
	return
}

func (s *{{$ModuleName}}) Create{{formatNameVal $ModuleName .}}(ptr *common.{{formatNameVal $ModuleName .}}Param, creater int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	{{formatPtrVal $ModuleName .}}Ptr := ptr.To{{formatNameVal $ModuleName .}}()
	{{formatPtrVal $ModuleName .}}Ptr.Creater = creater
	{{formatPtrVal $ModuleName .}}Ptr.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.{{toCamelCase $ModuleName}}Dao.Create{{formatNameVal $ModuleName .}}({{formatPtrVal $ModuleName .}}Ptr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{$ModuleName}}{{formatNameVal $ModuleName .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{$ModuleName}}) Update{{formatNameVal $ModuleName .}}(id int, ptr *common.{{formatNameVal $ModuleName .}}Param, updater int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	current{{formatNameVal $ModuleName .}}, currentErr := s.{{toCamelCase $ModuleName}}Dao.Query{{formatNameVal $ModuleName .}}(id, namespace)
	if currentErr != nil {
		return
	}

	current{{formatNameVal $ModuleName .}}.Creater = updater
	current{{formatNameVal $ModuleName .}}.UpdateTime = time.Now().UTC().Unix()

	ret, err = s.{{toCamelCase $ModuleName}}Dao.Update{{formatNameVal $ModuleName .}}(current{{formatNameVal $ModuleName .}}, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{$ModuleName}}{{formatNameVal $ModuleName .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{$ModuleName}}) Delete{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	ret, err = s.{{toCamelCase $ModuleName}}Dao.Delete{{formatNameVal $ModuleName .}}(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{$ModuleName}}{{formatNameVal $ModuleName .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
{{end}}`
	factory.RegisterFunc("formatNameVal", func(moduleName string, ptr *parser.Object) string {
		val := ptr.Name
		return val
	})
	factory.RegisterFunc("formatPtrVal", func(moduleName string, ptr *parser.Object) string {
		val := ptr.Name
		return util.ToCamelCase(val)
	})
	defer factory.UnRegisterFunc("formatPtrVal")
	defer factory.UnRegisterFunc("formatNameVal")

	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "Items": []*parser.Object{ptr}})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := util.ToTitle(ptr.FileName)
	commonFile := fmt.Sprintf("%s.go", "biz")
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(moduleName), "biz", commonFile)

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

		_, err = fileHandle.WriteString(s.declareBizImport(moduleName))
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

func (s *Biz) createObject(moduleName string, ptr *parser.Object) (err error) {
	commonFile := fmt.Sprintf("%s.go", util.ToLower(ptr.Name))
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(moduleName), "biz", commonFile)

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
