package util

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func DefineType(val *parser.Type) string {
	if val.IsArray {
		if val.IsPtr {
			return fmt.Sprintf("[]*model.%s", val.Name)
		}
		return fmt.Sprintf("[]model.%s", val.Name)
	}

	if val.IsPtr {
		return fmt.Sprintf("&model.%s", val.Name)
	}
	return fmt.Sprintf("model.%s", val.Name)
}

func DefineItem(val *parser.Item) string {
	//if val.IsPtr {
	//	return fmt.Sprintf("%s: &%v{}", val.Name, DefineType(val.Type))
	//}

	return fmt.Sprintf("%s: %v{}", val.Name, DefineType(val.Type))
}

func DefineObjectLite(ptr *parser.Object) (ret string) {
	objectDeclare :=
		`model.{{.Name}} { }`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"Name": ptr.Name, "Items": ptr.Items})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func DefineObject(ptr *parser.Object) (ret string) {
	objectDeclare :=
		`model.{{.Name}} { {{range .Items}}
{{- if isComplex .}} {{defineItem .}}, {{end -}}
{{end}} }`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"Name": ptr.Name, "Items": ptr.Items})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}
