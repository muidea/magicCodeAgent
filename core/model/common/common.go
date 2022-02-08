package common

import (
	"bytes"
	"fmt"
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
	moduleName := util.ToTitle(ptr.FileName)
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
	actionDeclare := `{{$ModuleName := .ModuleName}}
const (
{{range .Items}}
	Filter{{formatNameVal $ModuleName .}} = "/{{formatUrlVal $ModuleName .}}/query/"
	Query{{formatNameVal $ModuleName .}}  = "/{{formatUrlVal $ModuleName .}}/query/:id"
	Create{{formatNameVal $ModuleName .}} = "/{{formatUrlVal $ModuleName .}}/create/"
	Update{{formatNameVal $ModuleName .}} = "/{{formatUrlVal $ModuleName .}}/update/:id"
	Delete{{formatNameVal $ModuleName .}} = "/{{formatUrlVal $ModuleName .}}/delete/:id"
	Notify{{formatNameVal $ModuleName .}} = "/{{formatUrlVal $ModuleName .}}/notify/:id"
{{end}}
)`

	factory.RegisterFunc("formatNameVal", func(moduleName string, ptr *parser.Object) string {
		val := ptr.Name
		if moduleName != ptr.Name {
			val = fmt.Sprintf("%s%s", moduleName, ptr.Name)
		}

		return val
	})
	factory.RegisterFunc("formatUrlVal", func(moduleName string, ptr *parser.Object) string {
		val := ptr.Name
		if moduleName != ptr.Name {
			val = fmt.Sprintf("%s/%s", moduleName, ptr.Name)
		}

		return util.ToLower(val)
	})
	defer factory.UnRegisterFunc("formatUrlVal")
	defer factory.UnRegisterFunc("formatNameVal")

	moduleName := util.ToTitle(ptr.FileName)
	factoryPtr := factory.NewMemoryFactor(actionDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "Items": ptr.Objects})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Common) declareMessage(ptr *parser.Package) (ret string) {
	messageDeclare := `{{$BackQuote := .BackQuote}}{{$ModuleName := .ModuleName}}
{{range .Items}}

{{declareObjectView .}}

{{declareObjectViewAssign .}}

{{declareObjectLite .}}

{{declareObjectLiteAssign .}}

{{declareObjectParam .}}

{{declareObjectParamAssign .}}

type {{formatNameVal .}}Result struct {
	cd.Result
	{{formatNameVal .}} *{{formatNameVal .}}View {{$BackQuote}}json:"{{formatJsonVal . }}"{{$BackQuote}}
}

type {{formatNameVal .}}LiteListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{formatNameVal .}}  []*{{formatNameVal .}}Lite {{$BackQuote}}json:"{{formatJsonVal . }}"{{$BackQuote}}
}

type {{formatNameVal .}}ListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{formatNameVal .}}  []*{{formatNameVal .}}View {{$BackQuote}}json:"{{formatJsonVal . }}"{{$BackQuote}}
}

type {{formatNameVal .}}StatisticResult struct {
	{{formatNameVal .}}ListResult
	Summary []*TotalizerView {{$BackQuote}}json:"summary"{{$BackQuote}}
}
{{end}}
`
	factory.RegisterFunc("formatNameVal", func(ptr *parser.Object) string {
		val := ptr.Name
		return val
	})
	factory.RegisterFunc("formatJsonVal", func(ptr *parser.Object) string {
		val := ptr.Name
		return util.ToCamelCase(val)
	})
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
	defer factory.UnRegisterFunc("formatJsonVal")
	defer factory.UnRegisterFunc("formatNameVal")

	moduleName := util.ToTitle(ptr.FileName)
	factoryPtr := factory.NewMemoryFactor(messageDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleName": moduleName, "Items": ptr.Objects, "BackQuote": "`"})
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

	moduleName := ptr.FileName
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
