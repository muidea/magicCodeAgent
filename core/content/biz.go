package content

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
)

type Biz struct {
	projectName string
}

func (s *Biz) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$moduleList := .Define.Module}}
package biz

import (
	"github.com/muidea/{{$projectName}}/core/handler"
	{{range $mk, $mv := $moduleList}}
	{{$moduleName := $mv.Name}}
	{{$lowName := $mv.Lower $moduleName}}
	"github.com/muidea/{{$projectName}}/core/biz/{{$lowName}}"
	{{end}}
)

{{range $mk, $mv := $moduleList}}
{{$moduleName := $mv.Name}}
{{$lowName := $mv.Lower $moduleName}}
func New{{$moduleName}}(handler handler.{{$moduleName}}Handler) *{{$lowName}}.{{$moduleName}} {
	return {{$lowName}}.New{{$moduleName}}(handler)
}
{{end}}
`

	return
}

func (s *Biz) Create(define *Define, forceFlag bool) (err error) {
	if len(define.Module) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/biz")
	existFlag, existErr := util.Exists(filePath)
	if existErr != nil {
		err = existErr
		return
	}
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, "biz")
	existFlag, existErr = util.Exists(filePath)
	if existErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}

	if !existFlag || forceFlag {
		func() {
			fileHandle, fileErr := os.Create(filePath)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			templateHandle, templateErr := template.New("biz").Parse(s.getTemplate())
			if templateErr != nil {
				err = templateErr
				return
			}

			data := map[string]interface{}{"ProjectName": s.projectName, "Define": define}
			err = templateHandle.Execute(fileHandle, data)
		}()
		if err != nil {
			return
		}

		cmd := exec.Command("gofmt", "-w", filePath)
		err = cmd.Run()
	}

	return
}
