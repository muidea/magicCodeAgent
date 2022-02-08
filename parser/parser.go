package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strings"
)

type Type struct {
	Name    string `json:"name"`
	IsPtr   bool   `json:"isPtr"`
	IsArray bool   `json:"isArray"`
}

func (s *Type) String() string {
	val := fmt.Sprintf("%s", s.Name)
	if s.IsPtr {
		val = fmt.Sprintf("*%s", val)
	}
	if s.IsArray {
		val = fmt.Sprintf("[]%s", val)
	}

	return val
}

type Item struct {
	Name  string `json:"name"`
	Type  *Type  `json:"type"`
	IsPtr bool   `json:"isPtr"`
}

type Object struct {
	Name  string  `json:"name"`
	Items []*Item `json:"items"`
}

type Package struct {
	Name     string      `json:"packageName"`
	FileName string      `json:"fileName"`
	Objects  ObjectsList `json:"objects"`
}

type ObjectsList []*Object

func (s ObjectsList) Len() int {
	return len(s)
}

func (s ObjectsList) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s ObjectsList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func parseIdent(ptr *ast.Ident) (ret *Type, err error) {
	ret = &Type{Name: ptr.Name, IsPtr: false, IsArray: false}
	return
}

func parseArray(ptr *ast.ArrayType) (ret *Type, err error) {
	switch t := ptr.Elt.(type) {
	case *ast.Ident:
		ret, err = parseIdent(t)
		if err != nil {
			return
		}
		ret.IsArray = true
		return
	case *ast.StarExpr:
		ret, err = parseStart(t)
		if err != nil {
			return
		}
		if ret.IsArray {
			err = fmt.Errorf("illegal array, type:[]%v, pos:%v, end:%v", ret, t.Pos(), t.End())
			return
		}
		ret.IsArray = true
		return
	default:
		err = fmt.Errorf("illegal array, pos:%v, end:%v", t.Pos(), t.End())
	}

	err = fmt.Errorf("illegal array, pos:%v, end:%v", ptr.Pos(), ptr.End())
	return
}

func parseStart(ptr *ast.StarExpr) (ret *Type, err error) {
	switch t := ptr.X.(type) {
	case *ast.Ident:
		ret, err = parseIdent(t)
		if err != nil {
			return
		}
		ret.IsPtr = true
		return
	case *ast.ArrayType:
		ret, err = parseArray(t)
		if err != nil {
			return
		}
		if ret.IsPtr {
			err = fmt.Errorf("illegal start, type:*%v, pos:%v, end:%v", ret, t.Pos(), t.End())
			return
		}
		ret.IsPtr = true
		return
	default:
		err = fmt.Errorf("illegal start, pos:%v, end:%v", t.Pos(), t.End())
	}

	err = fmt.Errorf("illegal start, pos:%v, end:%v", ptr.Pos(), ptr.End())
	return
}

func parseField(ptr *ast.Field) (ret *Item, err error) {
	ret = &Item{}
	ret.Name = ptr.Names[0].Name

	switch t := ptr.Type.(type) {
	case *ast.Ident:
		ret.Type, err = parseIdent(t)
		return
	case *ast.ArrayType:
		ret.Type, err = parseArray(t)
		return
	case *ast.StarExpr:
		ret.Type, err = parseStart(t)
		return
	default:
		err = fmt.Errorf("illegal field define,name:%v, pos:%v, end:%v", ptr.Names[0].Name, ptr.Pos(), ptr.End())
	}

	err = fmt.Errorf("illegal field define,name:%v, pos:%v, end:%v", ptr.Names[0].Name, ptr.Pos(), ptr.End())
	return
}

func parseObject(ptr *ast.Object) (ret *Object, err error) {
	if ptr.Kind != ast.Typ {
		// skip invalid kind
		return
	}
	switch t := ptr.Decl.(type) {
	case *ast.TypeSpec:
		switch st := t.Type.(type) {
		case *ast.StructType:
			ret = &Object{}
			ret.Name = t.Name.Name
			for _, val := range st.Fields.List {
				itemVal, itemErr := parseField(val)
				if itemErr != nil {
					err = itemErr
					return
				}

				ret.Items = append(ret.Items, itemVal)
			}
		default:
			// skip,nothing
			return
		}
	default:
		return
	}

	return
}

func Parse(fullFileName string) (ret *Package, err error) {
	_, fileName := path.Split(fullFileName)
	strItems := strings.Split(fileName, ".")
	if len(strItems) != 2 || strItems[1] != "go" {
		err = fmt.Errorf("illegal file, fileName:%s", fullFileName)
		return
	}
	fileName = strItems[0]

	fSet := token.NewFileSet() // positions are relative to fSet
	astFile, astErr := parser.ParseFile(fSet, fullFileName, nil, parser.AllErrors)
	if astErr != nil {
		err = astErr
		return
	}

	ret = &Package{}
	ret.Name = astFile.Name.Name
	ret.FileName = fileName

	for _, val := range astFile.Scope.Objects {
		objPtr, objErr := parseObject(val)
		if objErr != nil {
			err = fmt.Errorf("illegal object, name:%s, reason:[%s]", objPtr.Name, objErr.Error())
			return
		}

		if objPtr != nil {
			ret.Objects = append(ret.Objects, objPtr)
		}
	}

	sort.Sort(ret.Objects)
	return
}
