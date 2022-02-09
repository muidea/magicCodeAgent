package common

import (
	"bytes"
	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

const (
	BackQuote = "`"
)

var reserved = []string{"Namespace"}
var entityItem = []string{"Creater", "Updater"}
var timeStamp = []string{"UpdateTime", "CreateTime"}
var specialItem = []string{"ID", "URL"}
var liteMaxIndex = 2

func isID(val string) bool {
	return val == "ID"
}

func isNumber(val string) bool {
	switch val {

	}

	return false
}

func isString(val string) bool {
	return val == "string"
}

func isBoolean(val string) bool {
	return val == "bool"
}

func isReserved(val string) bool {
	for _, sv := range reserved {
		if sv == val {
			return true
		}
	}

	return false
}

func isEntity(val string) bool {
	for _, sv := range entityItem {
		if sv == val {
			return true
		}
	}

	return false
}

func isSpecial(val string) bool {
	for _, sv := range specialItem {
		if sv == val {
			return true
		}
	}

	return false
}

func isTimeStamp(val string) bool {
	for _, sv := range timeStamp {
		if sv == val {
			return true
		}
	}

	return false
}

func DeclareItemView(val *parser.Item) string {
	arrayDeclare := `{{.NameVal}} []*{{.TypeVal}}Lite {{.BackQuote}}json:"{{.JsonVal}}"{{.BackQuote}}
`
	itemDeclare := `{{.NameVal}} *{{.TypeVal}}Lite {{.BackQuote}}json:"{{.JsonVal}}"{{.BackQuote}}
`
	simpleArrayDeclare := `{{.NameVal}} []{{.TypeVal}} {{.BackQuote}}json:"{{.JsonVal}}"{{.BackQuote}}
`
	simpleDeclare := `{{.NameVal}} {{.TypeVal}} {{.BackQuote}}json:"{{.JsonVal}}"{{.BackQuote}}
`
	jsonVal := util.ToCamelCase(val.Name)
	if isSpecial(val.Name) {
		jsonVal = util.ToLower(jsonVal)
	}
	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayDeclare)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name, "JsonVal": jsonVal, "BackQuote": BackQuote})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name, "JsonVal": jsonVal, "BackQuote": BackQuote})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayDeclare)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name, "JsonVal": jsonVal, "BackQuote": BackQuote})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name, "JsonVal": jsonVal, "BackQuote": BackQuote})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareEntityView(val *parser.Item) string {
	entityAssign := `{{.NameVal}}     *cc.EntityView {{.BackQuote}}json:"{{.JsonVal}}"{{.BackQuote}}
`
	jsonVal := util.ToCamelCase(val.Name)
	if isSpecial(val.Name) {
		jsonVal = util.ToLower(jsonVal)
	}
	factoryPtr := factory.NewMemoryFactor(entityAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name, "JsonVal": jsonVal, "BackQuote": "`"})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectView(ptr *parser.Object) (ret string) {
	viewDefine := `
type {{formatObjectName .ObjectVal}}View struct {
	{{.ContentVal -}}
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			byteBuff.WriteString(DeclareEntityView(val))
			continue
		}
		byteBuff.WriteString(DeclareItemView(val))
	}

	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLite(ptr *parser.Object) (ret string) {
	liteDefine := `
type {{formatObjectName .ObjectVal}}Lite struct {
	{{.ContentVal -}}
}
`
	byteBuff := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			byteBuff.WriteString(DeclareEntityView(val))
			continue
		}
		byteBuff.WriteString(DeclareItemView(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectParam(ptr *parser.Object) (ret string) {
	liteDefine := `
type {{formatObjectName .ObjectVal}}Param struct {
	{{.ContentVal -}}
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isTimeStamp(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			continue
		}

		byteBuff.WriteString(DeclareItemView(val))
	}
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemViewAssign(val *parser.Item) string {
	arrayAssign := `for _, val := range ptr.{{.NameVal}} {
		lite := &{{.TypeVal}}Lite{}
		lite.From{{.TypeVal}}(val)
		s.{{.NameVal}} = append(s.{{.NameVal}}, lite)
	}
`

	itemAssign := `s.{{.NameVal}} = &{{.TypeVal}}Lite{}
	s.{{.NameVal}}.From{{.TypeVal}}(ptr.{{.NameVal}})
`

	simpleAssign := `s.{{.NameVal}} = ptr.{{.NameVal}}
`
	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemLiteAssign(val *parser.Item) string {
	simpleAssign := `s.{{.NameVal}} = ptr.{{.NameVal}}
`
	factoryPtr := factory.NewMemoryFactor(simpleAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemLiteTo(val *parser.Item) string {
	arrayAssign := `for _, val := range s.{{.NameVal}} {
		ptr.{{.NameVal}} = append(ptr.{{.NameVal}}, val.To{{.TypeVal}}())
	}
`
	itemAssign := `if s.{{.NameVal}} != nil {
		ptr.{{.NameVal}} = s.{{.NameVal}}.To{{.TypeVal}}()
	}
`
	simpleNumAssign := `if s.{{.NameVal}} != 0 { 
		ptr.{{.NameVal}} = s.{{.NameVal}}
	}
`
	simpleStringAssign := `if s.{{.NameVal}} != "" { 
		ptr.{{.NameVal}} = s.{{.NameVal}}
	}
`
	simpleBooleanAssign := `if s.{{.NameVal}} { 
		ptr.{{.NameVal}} = s.{{.NameVal}}
	}
`
	simpleArrayAssign := `if len(s.{{.NameVal}}) > 0 { 
		ptr.{{.NameVal}} = s.{{.NameVal}}
	}
`

	if mu.IsComplex(val) {
		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if isNumber(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleNumAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if isString(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleStringAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if isBoolean(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleBooleanAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleNumAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareEntityViewAssign(val *parser.Item) string {
	entityAssign := `s.{{.NameVal}} = entityPtr
`
	factoryPtr := factory.NewMemoryFactor(entityAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectViewAssign(ptr *parser.Object) (ret string) {
	viewDefine := `
func (s *{{formatObjectName .ObjectVal}}View) From{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}, entityPtr *cc.EntityView) {
	{{.ContentVal -}}
	return
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			byteBuff.WriteString(DeclareEntityViewAssign(val))
			continue
		}
		byteBuff.WriteString(DeclareItemViewAssign(val))
	}

	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLiteAssign(ptr *parser.Object) (ret string) {
	liteDefine := `
func (s *{{formatObjectName .ObjectVal}}Lite) From{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}) {
	{{.ObjectAssignVal -}}
	return
}

func (s *{{formatObjectName .ObjectVal}}Lite) FromView(ptr *{{formatObjectName .ObjectVal}}View) {
	{{.ViewAssignVal -}}
	return
}

func (s *{{formatObjectName .ObjectVal}}Lite) To{{formatObjectName .ObjectVal}}() (ret *model.{{formatObjectName .ObjectVal}}) {
	ptr := &{{defineObject .ObjectVal}}
	{{.ObjectToVal -}}
	ret = ptr
	return
}
`
	objectAssign := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			continue
		}
		objectAssign.WriteString(DeclareItemViewAssign(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	viewAssign := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			continue
		}
		viewAssign.WriteString(DeclareItemLiteAssign(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	objectTo := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if isReserved(val.Name) {
			continue
		}
		if isTimeStamp(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			continue
		}

		objectTo.WriteString(DeclareItemLiteTo(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{
		"ObjectVal":       ptr,
		"ObjectAssignVal": objectAssign.String(),
		"ViewAssignVal":   viewAssign.String(),
		"ObjectToVal":     objectTo.String(),
	})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectParamAssign(ptr *parser.Object) (ret string) {
	paramToDefine := `
func (s *{{formatObjectName .ObjectVal}}Param) To{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}) (ret *model.{{formatObjectName .ObjectVal}}) {
	if ptr == nil {
		ptr = &{{defineObject .ObjectVal}}
	}

	{{.ContentVal}}
	ret = ptr
	return
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if isID(val.Name) {
			continue
		}

		if isReserved(val.Name) {
			continue
		}
		if isTimeStamp(val.Name) {
			continue
		}
		if isEntity(val.Name) {
			continue
		}

		byteBuff.WriteString(DeclareItemLiteTo(val))
	}

	factoryPtr := factory.NewMemoryFactor(paramToDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String()})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}
