package util

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func IsSimple(val *parser.Item) bool {
	switch val.Type.Name {
	case "bool",
		"int8", "int16", "int32", "int", "int64",
		"uint8", "uint16", "uint32", "uint", "uint64",
		"float32", "float64",
		"string", "time.Time":
		return true
	default:
	}

	return false
}

func DefineItem(val *parser.Item) string {
	if val.IsPtr {
		return fmt.Sprintf("%s: *model.%v{}", val.Name, val.Type)
	}

	return fmt.Sprintf("%s: model.%v{}", val.Name, val.Type)
}

func DefineObject(ptr *parser.Object) (ret string) {
	objectDeclare :=
		`model.{{.Name}} { {{range .Items}}
{{if !isSimple .}} {{defineItem .}}, {{end}}
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
