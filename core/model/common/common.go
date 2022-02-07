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

func (s *Common) declarePackage(ptr *parser.Package) (ret string) {
	packageDeclare := `
package common
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Common) declareImport(ptr *parser.Package) (ret string) {
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
	actionDeclare := `
const (
{{range .ItemsList}}
	Filter{{.NameVal}} = "/{{.UrlVal}}/query/"
	Query{{.NameVal}}  = "/{{.UrlVal}}/query/:id"
	Create{{.NameVal}} = "/{{.UrlVal}}/create/"
	Update{{.NameVal}} = "/{{.UrlVal}}/update/:id"
	Delete{{.NameVal}} = "/{{.UrlVal}}/delete/:id"
	Notify{{.NameVal}} = "/{{.UrlVal}}/notify/:id"
{{end}}
)`
	type Item struct {
		NameVal string
		UrlVal  string
	}

	itemList := []*Item{}
	moduleName := util.ToTitle(ptr.FileName)
	for _, val := range ptr.Objects {
		nameVal := val.Name
		if moduleName != val.Name {
			nameVal = fmt.Sprintf("%s%s", moduleName, val.Name)
		}

		urlVal := val.Name
		if moduleName != val.Name {
			urlVal = fmt.Sprintf("%s/%s", moduleName, val.Name)
		}
		urlVal = util.ToLower(urlVal)

		itemList = append(itemList, &Item{nameVal, urlVal})
	}

	factoryPtr := factory.NewMemoryFactor(actionDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemsList": itemList})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Common) declareMessage(ptr *parser.Package) (ret string) {
	messageDeclare := `{{$BackQuote := .BackQuote}}
{{range .ItemsList}}
type {{.NameVal}}View struct {
// TODO
}

func (s *{{.NameVal}}View) From{{.NameVal}}(ptr *model.{{.NameVal}}) {
// TODO
	return
}

type {{.NameVal}}Lite struct {
// TODO
}

func (s *{{.NameVal}}Lite) From{{.NameVal}}(ptr *model.{{.NameVal}}) {
// TODO
	return
}

func (s *{{.NameVal}}Lite) FromView(ptr *{{.NameVal}}View) {
// TODO
	return
}

func (s *{{.NameVal}}Lite) To{{.NameVal}}() (ret *model.{{.NameVal}}) {
// TODO
	return
}

type {{.NameVal}}Param struct {
// TODO
}

func (s *{{.NameVal}}Param) To{{.NameVal}}() (ret *model.{{.NameVal}}) {
// TODO
	return
}

type {{.NameVal}}Result struct {
	cd.Result
	{{.NameVal}} *{{.NameVal}}View {{$BackQuote}}json:"{{ toLower .NameVal }}"{{$BackQuote}}
}

type {{.NameVal}}LiteListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{.NameVal}}  []*{{.NameVal}}Lite {{$BackQuote}}json:"{{ toLower .NameVal }}"{{$BackQuote}}
}

type {{.NameVal}}ListResult struct {
	cd.Result
	Total int64       {{$BackQuote}}json:"total"{{$BackQuote}}
	{{.NameVal}}  []*{{.NameVal}}View {{$BackQuote}}json:"{{ toLower .NameVal }}"{{$BackQuote}}
}

{{end}}
`
	type Item struct {
		NameVal string
	}

	itemList := []*Item{}
	for _, val := range ptr.Objects {
		nameVal := val.Name
		itemList = append(itemList, &Item{nameVal})
	}

	factoryPtr := factory.NewMemoryFactor(messageDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemsList": itemList, "BackQuote": "`"})
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

		_, err = fileHandle.WriteString(s.declarePackage(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareImport(ptr))
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
