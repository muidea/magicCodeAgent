package data

import (
	"bytes"

	"github.com/muidea/magicCommon/foundation/log"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func getTypeValue(typeName string) (ret string) {
	switch typeName {
	case "int8", "int16", "int32", "int64", "int",
		"uint8", "uint16", "uint32", "uint64", "uint",
		"float32", "float64":
		ret = "number"
	case "bool":
		return "boolean"
	case "string":
		ret = "string"
	default:
		ret = "any"
	}

	return
}

func DeclareItem(val *parser.Item) string {
	arrayDeclare := `  {{.JsonVal}}?: any[];
`
	itemDeclare := `  {{.JsonVal}}?: any;
`
	simpleArrayDeclare := `  {{.JsonVal}}: {{.TypeVal}}[];
`
	simpleDeclare := `  {{.JsonVal}}: {{.TypeVal}};
`
	jsonVal := util.ToCamelCase(val.Name)
	if mu.IsSpecial(val.Name) {
		jsonVal = util.ToLower(jsonVal)
	}
	typeVal := getTypeValue(val.Type.Name)
	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayDeclare)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"TypeVal": typeVal, "JsonVal": jsonVal})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"TypeVal": typeVal, "JsonVal": jsonVal})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"TypeVal": typeVal, "JsonVal": jsonVal})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"TypeVal": typeVal, "JsonVal": jsonVal})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareEntity(val *parser.Item) string {
	entityAssign := `  {{.JsonVal}}?: any;
`
	jsonVal := util.ToCamelCase(val.Name)
	if mu.IsSpecial(val.Name) {
		jsonVal = util.ToLower(jsonVal)
	}
	factoryPtr := factory.NewMemoryFactor(entityAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"JsonVal": jsonVal})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectInterface(ptr *parser.Object) (ret string) {
	viewDefine := `export interface {{formatObjectName .ObjectVal}} {
{{  .ContentVal }}}`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
			byteBuff.WriteString(DeclareEntity(val))
			continue
		}

		byteBuff.WriteString(DeclareItem(val))
	}

	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}
