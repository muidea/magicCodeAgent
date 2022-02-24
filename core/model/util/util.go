package util

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

var reserved = []string{"Namespace"}
var entityItem = []string{"Creater", "Updater"}
var timeStamp = []string{"UpdateTime", "CreateTime", "DateTime"}
var specialItem = []string{"ID", "URL", "SKU", "SN", "WeChat"}
var innerFileName = []string{
	"common.go",
	"log.go",
	"totalizer.go",
	"helper.go",
}

func IsID(val string) bool {
	return val == "ID"
}

func IsNumber(val string) bool {
	switch val {

	}

	return false
}

func IsString(val string) bool {
	return val == "string"
}

func IsBoolean(val string) bool {
	return val == "bool"
}

func IsReserved(val string) bool {
	for _, sv := range reserved {
		if sv == val {
			return true
		}
	}

	return false
}

func IsEntity(val string) bool {
	for _, sv := range entityItem {
		if sv == val {
			return true
		}
	}

	return false
}

func IsSpecial(val string) bool {
	for _, sv := range specialItem {
		if sv == val {
			return true
		}
	}

	return false
}

func IsTimeStamp(val string) bool {
	for _, sv := range timeStamp {
		if sv == val {
			return true
		}
	}

	return false
}

func HasEntity(ptr *parser.Object) (ret bool) {
	ret = false
	for _, val := range ptr.Items {
		if IsEntity(val.Name) {
			ret = true
			return
		}
	}

	return
}

func IsInnerFile(fileName string) bool {
	for _, val := range innerFileName {
		if val == fileName {
			return true
		}
	}

	return false
}

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

func FormatModuleName(ptr *parser.Package) (ret string) {
	ret = util.ToTitle(ptr.FileName)
	return ret
}

func FormatModuleVal(ptr *parser.Package) (ret string) {
	ret = util.ToCamelCase(ptr.FileName)
	return ret
}

func FormatModulePath(ptr *parser.Package) (ret string) {
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
