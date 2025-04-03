package common

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/parser"
	"github.com/muidea/magicToolkit/pkg/factory"
)

type Common struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Common {
	return &Common{rootPath: rootPath, projectName: projectName}
}

func (s *Common) declarePackage() (ret string) {
	packageDeclare := `
package common
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Common) declareImport() (ret string) {
	importDeclare := `
import (
	cc "github.com/muidea/magicCas/common"
	cd "github.com/muidea/magicCommon/def"
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

func (s *Common) declareModule(ptr *parser.Package) (ret string) {
	moduleDeclare := `
const {{formatModuleName .ModuleVal}}Module = "/module/{{formatModulePath .ModuleVal}}"
`
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal

	return
}

func (s *Common) declareAction(ptr *parser.Package) (ret string) {
	actionDeclare := `{{$ModuleVal := .ModuleVal}}
const ({{with .ModuleVal}}
{{- range .Objects}}
	Filter{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionValue $ModuleVal .}}/query/"
	Query{{formatObjectActionName $ModuleVal .}}  = "/{{formatObjectActionValue $ModuleVal .}}/query/:id"
	Create{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionValue $ModuleVal .}}/create/"
	Update{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionValue $ModuleVal .}}/update/:id"
	Delete{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionValue $ModuleVal .}}/delete/:id"
	Notify{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionValue $ModuleVal .}}/notify/:id"
{{end}}{{end}})`
	factoryPtr := factory.NewMemoryFactor(actionDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Common) declareMessage(ptr *parser.Package) (ret string) {
	messageDeclare := `{{$BackQuote := .BackQuote}}{{$ModuleVal := .ModuleVal}}
{{- with $ModuleVal}}
{{- range .Objects }}
{{- if enableObjectView .}}
{{- declareObjectView .}}
{{- declareObjectViewAssign .}}
{{- end}}
{{- if enableObjectLite .}}
{{- declareObjectLite .}}
{{- declareObjectLiteAssign .}}
{{- if enableObjectView .}}
{{- declareObjectLiteFrom .}}
{{- end}}
{{- declareObjectLiteTo .}}
{{- end}}
{{- if enableObjectParam .}}
{{- declareObjectParam .}}
{{- declareObjectParamAssign .}}
{{- end}}
{{- if enableObjectView .}}
type {{formatObjectName .}}Result struct {
	cd.Result
	{{formatObjectName .}} *{{formatObjectName .}}View {{$BackQuote}}json:"{{formatObjectJson .}}"{{$BackQuote}}
}

type {{formatObjectName .}}ListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{formatObjectName .}}  []*{{formatObjectName .}}View {{$BackQuote}}json:"{{formatObjectJson .}}"{{$BackQuote}}
}

type {{formatObjectName .}}StatisticResult struct {
	{{formatObjectName .}}ListResult
	Summary []*TotalizerView {{$BackQuote}}json:"summary"{{$BackQuote}}
}
{{end}}
{{- if enableObjectLite .}}
type {{formatObjectName .}}LiteListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{formatObjectName .}}  []*{{formatObjectName .}}Lite {{$BackQuote}}json:"{{formatObjectJson .}}"{{$BackQuote}}
}
{{end}}
{{- end}}
{{- end}}
`
	factory.RegisterFunc("declareItemAssign", DeclareItemAssign)
	factory.RegisterFunc("declareItemFrom", DeclareItemFrom)
	factory.RegisterFunc("declareItemTo", DeclareItemTo)
	factory.RegisterFunc("declareObjectView", DeclareObjectView)
	factory.RegisterFunc("declareObjectViewAssign", DeclareObjectViewAssign)
	factory.RegisterFunc("declareObjectLite", DeclareObjectLite)
	factory.RegisterFunc("declareObjectLiteAssign", DeclareObjectLiteAssign)
	factory.RegisterFunc("declareObjectLiteFrom", DeclareObjectLiteFrom)
	factory.RegisterFunc("declareObjectLiteTo", DeclareObjectLiteTo)
	factory.RegisterFunc("declareObjectParam", DeclareObjectParam)
	factory.RegisterFunc("declareObjectParamAssign", DeclareObjectParamTo)

	defer factory.UnRegisterFunc("declareObjectParamAssign")
	defer factory.UnRegisterFunc("declareObjectParam")
	defer factory.UnRegisterFunc("declareObjectLiteTo")
	defer factory.UnRegisterFunc("declareObjectLiteFrom")
	defer factory.UnRegisterFunc("declareObjectLiteAssign")
	defer factory.UnRegisterFunc("declareObjectLite")
	defer factory.UnRegisterFunc("declareObjectView")
	defer factory.UnRegisterFunc("declareObjectViewAssign")
	defer factory.UnRegisterFunc("declareItemTo")
	defer factory.UnRegisterFunc("declareItemFrom")
	defer factory.UnRegisterFunc("declareItemAssign")

	factoryPtr := factory.NewMemoryFactor(messageDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "BackQuote": "`"})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Common) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := mu.FormatModulePath(ptr)
	commonFile := fmt.Sprintf("%s.go", moduleName)
	commonFilePath := path.Join(s.rootPath, "common", commonFile)
	if !util.Exists(commonFilePath) {
		err = fmt.Errorf("%s not exists", commonFile)
		return
	}

	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
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

		_, err = fileHandle.WriteString(s.declareModule(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareAction(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareMessage(ptr))
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
