package common

import (
	log "github.com/cihub/seelog"

	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func DeclareObjectView(ptr *parser.Object) (ret string) {
	viewDefine := `{{$ObjectVal := .ObjectVal}}{{with $ObjectVal}}
type {{formatObjectName $ObjectVal}}View struct {
{{- range .Items}}
{{- if enableView .}}
	{{declareItem .}}
{{- end}}
{{- end}}
}{{end}}
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
{{- range .Items}}
{{- if enableLite .}}
	{{declareItem .}}
{{- end}}
{{- end}}
}{{end}}
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
{{- range .Items}}
{{- if enableParam .}}
	{{declareItem .}}
{{- end}}
{{- end}}
}{{end}}
`
	factoryPtr := factory.NewMemoryFactor(paramDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemAssign(val *parser.Item) string {
	arrayPtrAssign := `for _, val := range ptr.{{formatItemName .ItemVal}} {
		lite := &{{formatItemType .ItemVal}}Lite{}
		lite.From{{formatItemType .ItemVal}}(val)
		s.{{formatItemName .ItemVal}} = append(s.{{formatItemName .ItemVal}}, lite)
	}`
	arrayAssign := `for _, val := range ptr.{{formatItemName .ItemVal}} {
		lite := &{{formatItemType .ItemVal}}Lite{}
		lite.From{{formatItemType .ItemVal}}(&val)
		s.{{formatItemName .ItemVal}} = append(s.{{formatItemName .ItemVal}}, lite)
	}`
	itemPtrAssign := `if ptr.{{formatItemName .ItemVal}} != nil {
		s.{{formatItemName .ItemVal}} = &{{formatItemType .ItemVal}}Lite{}
		s.{{formatItemName .ItemVal}}.From{{formatItemType .ItemVal}}(ptr.{{formatItemName .ItemVal}})
	}`
	itemAssign := `s.{{formatItemName .ItemVal}} = &{{formatItemType .ItemVal}}Lite{}
	s.{{formatItemName .ItemVal}}.From{{formatItemType .ItemVal}}(&ptr.{{formatItemName .ItemVal}})`
	simpleAssign := `s.{{formatItemName .ItemVal}} = ptr.{{formatItemName .ItemVal}}`

	if mu.IsComplex(val) {
		if val.Type.IsPtr {
			if val.Type.IsArray {
				factoryPtr := factory.NewMemoryFactor(arrayPtrAssign)
				contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
				if contentErr != nil {
					log.Errorf("create content failed, err:%s", contentErr.Error())
					return ""
				}

				return contentVal
			}

			factoryPtr := factory.NewMemoryFactor(itemPtrAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemFrom(val *parser.Item) string {
	entityAssign := `s.{{formatItemName .ItemVal}} = ptr.{{formatItemName .ItemVal}}`
	factoryPtr := factory.NewMemoryFactor(entityAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareItemTo(val *parser.Item) string {
	arrayPtrAssign := `if len(s.{{formatItemName .ItemVal}}) > 0 {
		ptr.{{formatItemName .ItemVal}} = []*model.{{formatItemType .ItemVal}}{}
		for _, val := range s.{{formatItemName .ItemVal}} {
			ptr.{{formatItemName .ItemVal}} = append(ptr.{{formatItemName .ItemVal}}, val.To{{formatItemType .ItemVal}}())
		}
	}`
	arrayAssign := `if len(s.{{formatItemName .ItemVal}}) > 0 {
		ptr.{{formatItemName .ItemVal}} = []model.{{formatItemType .ItemVal}}{}
		for _, val := range s.{{formatItemName .ItemVal}} {
			ptr.{{formatItemName .ItemVal}} = append(ptr.{{formatItemName .ItemVal}}, *val.To{{formatItemType .ItemVal}}())
		}
	}`
	itemPtrAssign := `if s.{{formatItemName .ItemVal}} != nil {
		ptr.{{formatItemName .ItemVal}} = s.{{formatItemName .ItemVal}}.To{{formatItemType .ItemVal}}()
	}`
	itemAssign := `if s.{{formatItemName .ItemVal}} != nil {
		ptr.{{formatItemName .ItemVal}} = *s.{{formatItemName .ItemVal}}.To{{formatItemType .ItemVal}}()
	}`
	simpleNumAssign := `if s.{{formatItemName .ItemVal}} != 0 { 
		ptr.{{formatItemName .ItemVal}} = s.{{formatItemName .ItemVal}}
	}`
	simpleStringAssign := `if s.{{formatItemName .ItemVal}} != "" { 
		ptr.{{formatItemName .ItemVal}} = s.{{formatItemName .ItemVal}}
	}`
	simpleBooleanAssign := `if s.{{formatItemName .ItemVal}} { 
		ptr.{{formatItemName .ItemVal}} = s.{{formatItemName .ItemVal}}
	}`
	simpleArrayAssign := `if len(s.{{formatItemName .ItemVal}}) > 0 { 
		ptr.{{formatItemName .ItemVal}} = s.{{formatItemName .ItemVal}}
	}`

	if mu.IsComplex(val) {
		if val.Type.IsPtr {
			if val.Type.IsArray {
				factoryPtr := factory.NewMemoryFactor(arrayPtrAssign)
				contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
				if contentErr != nil {
					log.Errorf("create content failed, err:%s", contentErr.Error())
					return ""
				}

				return contentVal
			}

			factoryPtr := factory.NewMemoryFactor(itemPtrAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		if val.Type.IsArray {
			factoryPtr := factory.NewMemoryFactor(arrayAssign)
			contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
			if contentErr != nil {
				log.Errorf("create content failed, err:%s", contentErr.Error())
				return ""
			}

			return contentVal
		}

		factoryPtr := factory.NewMemoryFactor(itemAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if val.Type.IsArray {
		factoryPtr := factory.NewMemoryFactor(simpleArrayAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if mu.IsNumber(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleNumAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}
	if mu.IsString(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleStringAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	if mu.IsBoolean(val.Type.Name) {
		factoryPtr := factory.NewMemoryFactor(simpleBooleanAssign)
		contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
		if contentErr != nil {
			log.Errorf("create content failed, err:%s", contentErr.Error())
			return ""
		}

		return contentVal
	}

	factoryPtr := factory.NewMemoryFactor(simpleNumAssign)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": val})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectViewAssign(ptr *parser.Object) (ret string) {
	viewDefine := `{{$ObjectVal := .ObjectVal}}
func (s *{{formatObjectName $ObjectVal}}View) From{{formatObjectName $ObjectVal}}(ptr *model.{{formatObjectName $ObjectVal}}) {
	if ptr == nil {
		return
	}
{{- with $ObjectVal}}
{{range .Items}}
{{- if enableView .}}
{{- if isNormalItem .}}
	{{declareItemAssign .}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
	return
}
`
	factoryPtr := factory.NewMemoryFactor(viewDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLiteAssign(ptr *parser.Object) (ret string) {
	liteDefine := `{{$ObjectVal := .ObjectVal}}
func (s *{{formatObjectName $ObjectVal}}Lite) From{{formatObjectName $ObjectVal}}(ptr *model.{{formatObjectName $ObjectVal}}) {
	if ptr == nil {
		return
	}
{{- with $ObjectVal}}
{{range .Items}}
{{- if enableLite .}}
{{- if isNormalItem .}}
	{{declareItemAssign .}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
	return
}
`
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLiteFrom(ptr *parser.Object) (ret string) {
	liteDefine := `{{$ObjectVal := .ObjectVal}}
func (s *{{formatObjectName $ObjectVal}}Lite) FromView(ptr *{{formatObjectName $ObjectVal}}View) {
	if ptr == nil {
		return
	}
{{- with $ObjectVal}}
{{range .Items}}
{{- if enableLite .}}
{{- if isNormalItem .}}
	{{declareItemFrom .}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
	return
}
`
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectLiteTo(ptr *parser.Object) (ret string) {
	liteDefine := `{{$ObjectVal := .ObjectVal}}
func (s *{{formatObjectName $ObjectVal}}Lite) To{{formatObjectName $ObjectVal}}() (ret *model.{{formatObjectName $ObjectVal}}) {
	ptr := &model.{{formatObjectName $ObjectVal}}{}
{{- with $ObjectVal -}}
{{range .Items}}
{{- if enableLite .}}
{{- if isNormalItem .}}
	{{declareItemTo .}}
{{- end}}
{{- end}}
{{- end}}
{{end}}
	ret = ptr
	return
}
`
	factoryPtr := factory.NewMemoryFactor(liteDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}

func DeclareObjectParamTo(ptr *parser.Object) (ret string) {
	paramToDefine := `{{$ObjectVal := .ObjectVal}}
func (s *{{formatObjectName .ObjectVal}}Param) To{{formatObjectName .ObjectVal}}(ptr *model.{{formatObjectName .ObjectVal}}) (ret *model.{{formatObjectName .ObjectVal}}) {
	if ptr == nil {
		ptr = &model.{{formatObjectName $ObjectVal}}{}
	}
{{- with $ObjectVal -}}
{{range .Items}}
{{- if enableParam .}}
{{- if isNormalItem .}}
	{{declareItemTo .}}
{{- end}}
{{- end}}
{{- end}}
{{end}}
	ret = ptr
	return
}
`
	factoryPtr := factory.NewMemoryFactor(paramToDefine)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ObjectVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return ""
	}

	return contentVal
}
