package dao

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

type Dao struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Dao {
	return &Dao{rootPath: rootPath, projectName: projectName}
}

func (s *Dao) declarePackage() (ret string) {
	packageDeclare := `
package dao
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Dao) declareImport() (ret string) {
	importDeclare := `
import (
	"fmt"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	"github.com/muidea/{{.ProjectName}}/model"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) declareInterface(ptr *parser.Package) (ret string) {
	interfaceDeclare := `{{$ModuleName:=.ModuleName}}
type {{$ModuleName}} interface {
{{range .Items}}
	Filter{{formatNameVal $ModuleName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatNameVal $ModuleName .}}, total int64, err error)
	Query{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error)
	Create{{formatNameVal $ModuleName .}}(ptr *model.{{formatNameVal $ModuleName .}}, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error)
	Update{{formatNameVal $ModuleName .}}(ptr *model.{{formatNameVal $ModuleName .}}, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error)
	Delete{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error)
{{end}}
}`
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

	moduleName := util.ToTitle(ptr.FileName)
	factoryPtr := factory.NewMemoryFactor(interfaceDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "Items": ptr.Objects})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `
func New(clnt client.Client) {{.ModuleName}} {
	return &{{toCamelCase .ModuleName}}{batisClient: clnt}
}

type {{toCamelCase .ModuleName}} struct {
	batisClient client.Client
}`
	type Item struct {
		NameVal string
	}

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

func (s *Dao) declareFunc(ptr *parser.Package) (ret string) {
	funcDeclare := `{{$ModuleName := .ModuleName}}
{{range .Items}}
func (s *{{toCamelCase $ModuleName}}) Filter{{formatNameVal $ModuleName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatNameVal $ModuleName .}}, total int64, err error) {
	ret = []*model.{{formatNameVal $ModuleName .}}{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *{{toCamelCase $ModuleName}}) Query{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &model.{{formatNameVal $ModuleName .}}{ID: id, Namespace: namespace}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{toCamelCase $ModuleName}}) Create{{formatNameVal $ModuleName .}}(ptr *model.{{formatNameVal $ModuleName .}}, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{toCamelCase $ModuleName}}) Update{{formatNameVal $ModuleName .}}(ptr *model.{{formatNameVal $ModuleName .}}, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{toCamelCase $ModuleName}}) Delete{{formatNameVal $ModuleName .}}(id int, namespace string) (ret *model.{{formatNameVal $ModuleName .}}, err error) {
	valPtr, valErr := s.Query{{formatNameVal $ModuleName .}}(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
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

	moduleName := util.ToTitle(ptr.FileName)
	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "Items": ptr.Objects})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := ptr.FileName
	commonFile := fmt.Sprintf("%s.go", "dao")
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(moduleName), "dao", commonFile)

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

		_, err = fileHandle.WriteString(s.declareImport())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareInterface(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareDefine(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareFunc(ptr))
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
