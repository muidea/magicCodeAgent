package util

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func IsComplex(val *parser.Item) bool {
	switch val.Type.Name {
	case "bool",
		"int8", "int16", "int32", "int", "int64",
		"uint8", "uint16", "uint32", "uint", "uint64",
		"float32", "float64",
		"string", "time.Time":
		return false
	default:
	}

	return true
}

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
	if val.IsPtr {
		return fmt.Sprintf("%s: &%v{}", val.Name, DefineType(val.Type))
	}

	return fmt.Sprintf("%s: %v{}", val.Name, DefineType(val.Type))
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

func FormatModuleName(ptr *parser.Package) (ret string) {
	ret = util.ToTitle(ptr.FileName)
	return ret
}

func FormatModuleVal(ptr *parser.Package) (ret string) {
	ret = util.ToLower(ptr.FileName)
	return ret
}

func FormatObjectName(ptr *parser.Object) string {
	val := ptr.Name
	return val
}

func FormatObjectVal(ptr *parser.Object) string {
	val := util.ToCamelCase(ptr.Name)
	return val
}

func FormatObjectJson(ptr *parser.Object) string {
	val := util.ToCamelCase(FormatObjectName(ptr))
	return val
}

func FormatObjectActionName(modulePtr *parser.Package, ptr *parser.Object) (ret string) {
	val := ptr.Name
	moduleName := FormatModuleName(modulePtr)
	if moduleName != ptr.Name {
		val = fmt.Sprintf("%s%s", moduleName, ptr.Name)
	}

	return val
}

func FormatObjectActionVal(modulePtr *parser.Package, ptr *parser.Object) (ret string) {
	val := ptr.Name
	moduleName := FormatModuleName(modulePtr)
	if moduleName != ptr.Name {
		val = fmt.Sprintf("%s/%s", moduleName, ptr.Name)
	}

	return util.ToLower(val)
}
