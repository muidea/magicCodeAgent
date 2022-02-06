package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

type Item struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	IsPtr   bool   `json:"isPtr"`
	IsArray bool   `json:"isArray"`
}

type Object struct {
	Name  string  `json:"name"`
	Items []*Item `json:"items"`
}

type Package struct {
	Name     string    `json:"packageName"`
	FileName string    `json:"fileName"`
	Objects  []*Object `json:"objects"`
}

func parseIdent(ptr *ast.Ident) (ret string, err error) {
	ret = ptr.Name
	return
}

func parseArray(ptr *ast.ArrayType) (ret string, err error) {
	switch t := ptr.Elt.(type) {
	case *ast.Ident:
		ret = t.Name
		return
	case *ast.StarExpr:
		strVal, bVal, valErr := parseStart(t)
		if valErr != nil {
			err = valErr
			return
		}
		if !bVal {
			ret = strVal
			return
		}
		err = fmt.Errorf("illegal array element, pos:%v, end:%v", t.Pos(), t.End())
	default:
		err = fmt.Errorf("illegal array, pos:%v, end:%v", t.Pos(), t.End())
	}

	err = fmt.Errorf("illegal array, pos:%v, end:%v", ptr.Pos(), ptr.End())
	return
}

func parseStart(ptr *ast.StarExpr) (strVal string, bVal bool, err error) {
	switch t := ptr.X.(type) {
	case *ast.Ident:
		strVal = t.Name
		bVal = false
		return
	case *ast.ArrayType:
		strVal, err = parseArray(t)
		bVal = true
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
		ret.IsPtr = false
		ret.IsArray = false
		ret.Type, err = parseIdent(t)
		return
	case *ast.ArrayType:
		ret.IsPtr = false
		ret.IsArray = true
		ret.Type, err = parseArray(t)
		return
	case *ast.StarExpr:
		strVal, bVal, errVal := parseStart(t)
		if errVal != nil {
			err = errVal
			return
		}
		if bVal {
			err = fmt.Errorf("illegal field define,name:%v, pos:%v, end:%v", ptr.Names[0].Name, ptr.Pos(), ptr.End())
			return
		}

		ret.IsPtr = true
		ret.IsArray = bVal
		ret.Type = strVal
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

	return
}
