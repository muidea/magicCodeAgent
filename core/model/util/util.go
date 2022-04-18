package util

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/parser"
)

var reservedItem = []string{"Namespace"}
var entityItem = []string{"Creater", "Updater"}
var specialItem = []string{"ID", "URL", "SKU", "SN", "WeChat"}
var timeStampItem = []string{"UpdateTime", "CreateTime", "DateTime"}
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
	case "int8", "int16", "int32", "int64", "int",
		"uint8", "uint16", "uint32", "uint64", "uint",
		"float32", "float64":
		return true
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
	for _, sv := range reservedItem {
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
	for _, sv := range timeStampItem {
		if sv == val {
			return true
		}
	}

	return false
}

func EnableView(ptr *parser.Item) bool {
	vals, ok := ptr.Tags["tag"]
	if !ok {
		return false
	}

	for _, val := range vals {
		if val == "view" {
			return true
		}
	}

	return false
}

func EnableLite(ptr *parser.Item) bool {
	vals, ok := ptr.Tags["tag"]
	if !ok {
		return false
	}

	for _, val := range vals {
		if val == "lite" {
			return true
		}
	}

	return false
}

func EnableParam(ptr *parser.Item) bool {
	vals, ok := ptr.Tags["tag"]
	if !ok {
		return false
	}

	for _, val := range vals {
		if val == "param" {
			return true
		}
	}

	return false
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

func FormatModuleName(ptr *parser.Package) (ret string) {
	ret = util.ToTitle(ptr.FileName)
	return ret
}

func FormatModuleValue(ptr *parser.Package) (ret string) {
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

func FormatObjectValue(ptr *parser.Object) string {
	val := util.ToCamelCase(ptr.Name)
	return val
}

func FormatObjectJson(ptr *parser.Object) string {
	val := util.ToCamelCase(ptr.Name)
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

func FormatObjectActionValue(modulePtr *parser.Package, ptr *parser.Object) (ret string) {
	val := ptr.Name
	moduleName := FormatModuleName(modulePtr)
	if moduleName != ptr.Name {
		val = fmt.Sprintf("%s/%s", moduleName, ptr.Name)
	}

	return util.ToLower(val)
}

func FormatItemName(ptr *parser.Item) string {
	return ptr.Name
}

func FormatItemValue(ptr *parser.Item) string {
	val := util.ToCamelCase(ptr.Name)
	return val
}

func FormatItemJson(ptr *parser.Item) string {
	val := util.ToCamelCase(ptr.Name)
	if IsSpecial(ptr.Name) {
		val = util.ToLower(val)
	}
	return val
}
