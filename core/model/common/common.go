package common

import (
	"bytes"
	"fmt"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"os"
	"os/exec"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
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
const {{.ModuleName}}Module = "/module/{{toLower .ModuleName}}"
`
	moduleName := mu.FormatModuleName(ptr)
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName})
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
{{range .Objects}}
	Filter{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionVal $ModuleVal .}}/query/"
	Query{{formatObjectActionName $ModuleVal .}}  = "/{{formatObjectActionVal $ModuleVal .}}/query/:id"
	Create{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionVal $ModuleVal .}}/create/"
	Update{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionVal $ModuleVal .}}/update/:id"
	Delete{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionVal $ModuleVal .}}/delete/:id"
	Notify{{formatObjectActionName $ModuleVal .}} = "/{{formatObjectActionVal $ModuleVal .}}/notify/:id"
{{end}}{{end}}
)`
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
	messageDeclare := `{{$BackQuote := .BackQuote}}{{$ModuleVal := .ModuleVal}}{{with $ModuleVal}}
{{range .Objects}}

{{declareObjectView .}}

{{declareObjectViewAssign .}}

{{declareObjectLite .}}

{{declareObjectLiteAssign .}}

{{declareObjectParam .}}

{{declareObjectParamAssign .}}

type {{formatObjectName .}}Result struct {
	cd.Result
	{{formatObjectName .}} *{{formatObjectName .}}View {{$BackQuote}}json:"{{formatObjectJson .}}"{{$BackQuote}}
}

type {{formatObjectName .}}LiteListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{formatObjectName .}}  []*{{formatObjectName .}}Lite {{$BackQuote}}json:"{{formatObjectJson .}}"{{$BackQuote}}
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
{{end}}{{end}}
`
	factory.RegisterFunc("declareObjectView", DeclareObjectView)
	factory.RegisterFunc("declareObjectViewAssign", DeclareObjectViewAssign)
	factory.RegisterFunc("declareObjectLite", DeclareObjectLite)
	factory.RegisterFunc("declareObjectLiteAssign", DeclareObjectLiteAssign)
	factory.RegisterFunc("declareObjectParam", DeclareObjectParam)
	factory.RegisterFunc("declareObjectParamAssign", DeclareObjectParamAssign)

	defer factory.UnRegisterFunc("declareObjectParam")
	defer factory.UnRegisterFunc("declareObjectParamAssign")
	defer factory.UnRegisterFunc("declareObjectLite")
	defer factory.UnRegisterFunc("declareObjectLiteAssign")
	defer factory.UnRegisterFunc("declareObjectView")
	defer factory.UnRegisterFunc("declareObjectViewAssign")

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

	log.Info(commonFilePath)
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
