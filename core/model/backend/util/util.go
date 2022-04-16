package util

import (
	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

const (
	BackQuote = "`"
)

func IsEntityItem(val *parser.Item) bool {
	return mu.IsEntity(val.Type.Name)
}

func FormatItemType(val *parser.Item) string {
	return val.Type.Name
}

func DeclareItem(val *parser.Item) string {
	entityDeclare := `{{formatItemName .ItemVal}} *cc.EntityView {{.BackQuote}}json:"{{formatItemJson .ItemVal}}"{{.BackQuote}}`
	arrayDeclare := `{{formatItemName .ItemVal}} []*{{formatItemType .ItemVal}}Lite {{.BackQuote}}json:"{{formatItemJson .ItemVal}}"{{.BackQuote}}`
	itemDeclare := `{{formatItemName .ItemVal}} *{{formatItemType .ItemVal}}Lite {{.BackQuote}}json:"{{formatItemJson .ItemVal}}"{{.BackQuote}}`
	simpleArrayDeclare := `{{formatItemName .ItemVal}} []{{formatItemType .ItemVal}} {{.BackQuote}}json:"{{formatItemJson .ItemVal}}"{{.BackQuote}}`
	simpleDeclare := `{{formatItemName .ItemVal}} {{formatItemType .ItemVal}} {{.BackQuote}}json:"{{formatItemJson .ItemVal}}"{{.BackQuote}}`
	if IsEntityItem(val) {
		factoryPtr := factory.NewMemoryFactor(entityDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val, "BackQuote": BackQuote})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayDeclare)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val, "BackQuote": BackQuote})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val, "BackQuote": BackQuote})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val, "BackQuote": BackQuote})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val, "BackQuote": BackQuote})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DefineItem(val *parser.Item) string {
	arrayDeclare := `{{formatItemValue .ItemVal}} :=[]*{{formatItemType .ItemVal}}`
	itemDeclare := `{{formatItemValue .ItemVal}} := &{{formatItemType .ItemVal}}`
	simpleArrayDeclare := `{{formatItemValue .ItemVal}} := []{{formatItemType .ItemVal}}`
	simpleDeclare := `{{formatItemValue .ItemVal}} := {{formatItemType .ItemVal}}`
	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayDeclare)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
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
