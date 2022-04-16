package common

import (
	"bytes"
	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

var liteMaxIndex = 2

func DeclareObjectView(ptr *parser.Object) (ret string) {
	viewDefine := `{{$ObjectVal := .ObjectVal}}{{with $ObjectVal}}
type {{formatObjectName $ObjectVal}}View struct {
{{range .Items}}
	{{declareItem .}}
{{end}}
}
{{end}}
`
	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLite(ptr *parser.Object) (ret string) {
	liteDefine := `{{$ObjectVal := .ObjectVal}}{{with $ObjectVal}}
type {{formatObjectName $ObjectVal}}Lite struct {
{{range .Items}}
	{{declareItem .}}
{{end}}
}
{{end}}
`
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectParam(ptr *parser.Object) (ret string) {
	paramDefine := `{{$ObjectVal := .ObjectVal}}{{with $ObjectVal}}
type {{formatObjectName $ObjectVal}}Param struct {
{{range .Items}}
	{{declareItem .}}
{{end}}
}
{{end}}
`
	factoryPtr := factory.NewMemoryFactor(paramDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemViewAssign(val *parser.Item) string {
	arrayPtrAssign := `for _, val := range ptr.{{.NameVal}} {
		lite := &{{.TypeVal}}Lite{}
		lite.From{{.TypeVal}}(val)
		s.{{.NameVal}} = append(s.{{.NameVal}}, lite)
	}
`
	arrayAssign := `for _, val := range ptr.{{.NameVal}} {
		lite := &{{.TypeVal}}Lite{}
		lite.From{{.TypeVal}}(&val)
		s.{{.NameVal}} = append(s.{{.NameVal}}, lite)
	}
`
	itemPtrAssign := `if ptr.{{.NameVal}} != nil {
		s.{{.NameVal}} = &{{.TypeVal}}Lite{}
		s.{{.NameVal}}.From{{.TypeVal}}(ptr.{{.NameVal}})
	}
`
	itemAssign := `s.{{.NameVal}} = &{{.TypeVal}}Lite{}
	s.{{.NameVal}}.From{{.TypeVal}}(&ptr.{{.NameVal}})
`
	simpleAssign := `s.{{.NameVal}} = ptr.{{.NameVal}}
`
	if mu.IsComplex(val) {
		if val.Type.IsPtr {
			if val.Type.IsArray {
				factoryPtr := factory.NewMemoryFactor(arrayPtrAssign)
				contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
				if contentErr != nil {
					log.Errorf("create content failed, err:%s", contentErr.Error())
					return ""
				}

				return contentVal
			}

			factoryPtr := factory.NewMemoryFactor(itemPtrAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

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
	arrayPtrAssign := `if len(s.{{.NameVal}}) > 0 {
		ptr.{{.NameVal}} = []*model.{{.TypeVal}}{}
		for _, val := range s.{{.NameVal}} {
			ptr.{{.NameVal}} = append(ptr.{{.NameVal}}, val.To{{.TypeVal}}())
		}
	}
`
	arrayAssign := `if len(s.{{.NameVal}}) > 0 {
		ptr.{{.NameVal}} = []model.{{.TypeVal}}{}
		for _, val := range s.{{.NameVal}} {
			ptr.{{.NameVal}} = append(ptr.{{.NameVal}}, *val.To{{.TypeVal}}())
		}
	}
`
	itemPtrAssign := `if s.{{.NameVal}} != nil {
		ptr.{{.NameVal}} = s.{{.NameVal}}.To{{.TypeVal}}()
	}
`
	itemAssign := `if s.{{.NameVal}} != nil {
		ptr.{{.NameVal}} = *s.{{.NameVal}}.To{{.TypeVal}}()
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
		if val.Type.IsPtr {
			if val.Type.IsArray {
				factoryPtr := factory.NewMemoryFactor(arrayPtrAssign)
				contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
				if contentErr != nil {
					log.Errorf("create content failed, err:%s", contentErr.Error())
					return ""
				}

				return contentVal
			}

			factoryPtr := factory.NewMemoryFactor(itemPtrAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

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
	if mu.IsNumber(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleNumAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if mu.IsString(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleStringAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"NameVal": val.Name, "TypeVal": val.Type.Name})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if mu.IsBoolean(val.Type.Name) {
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
{{- if .HasEntity -}}
func (s *{{formatObjectName .ObjectVal}}View) From{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}, entityPtr *cc.EntityView) {
{{- else -}}
func (s *{{formatObjectName .ObjectVal}}View) From{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}) {
{{- end}}
	if ptr == nil {
		return
	}
	{{.ContentVal -}}
	return
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
			byteBuff.WriteString(DeclareEntityViewAssign(val))
			continue
		}
		byteBuff.WriteString(DeclareItemViewAssign(val))
	}

	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr, "ContentVal": byteBuff.String(), "HasEntity": mu.HasEntity(ptr)})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLiteAssign(ptr *parser.Object) (ret string) {
	liteDefine := `
func (s *{{formatObjectName .ObjectVal}}Lite) From{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}) {
	if ptr == nil {
		return
	}
	{{.ObjectAssignVal -}}
	return
}

func (s *{{formatObjectName .ObjectVal}}Lite) FromView(ptr *{{formatObjectName .ObjectVal}}View) {
	if ptr == nil {
		return
	}
	{{.ViewAssignVal -}}
	return
}

func (s *{{formatObjectName .ObjectVal}}Lite) To{{formatObjectName .ObjectVal}}() (ret *model.{{formatObjectName .ObjectVal}}) {
	ptr := &{{defineObjectLite .ObjectVal}}
	{{.ObjectToVal -}}
	ret = ptr
	return
}
`
	objectAssign := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
			continue
		}
		objectAssign.WriteString(DeclareItemViewAssign(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	viewAssign := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
			continue
		}
		viewAssign.WriteString(DeclareItemLiteAssign(val))

		if idx >= liteMaxIndex {
			break
		}
	}
	objectTo := bytes.Buffer{}
	for idx, val := range ptr.Items {
		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsTimeStamp(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
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
		ptr = &{{defineObjectLite .ObjectVal}}
	}

	{{.ContentVal}}
	ret = ptr
	return
}
`
	byteBuff := bytes.Buffer{}
	for _, val := range ptr.Items {
		if mu.IsID(val.Name) {
			continue
		}

		if mu.IsReserved(val.Name) {
			continue
		}
		if mu.IsTimeStamp(val.Name) {
			continue
		}
		if mu.IsEntity(val.Name) {
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
