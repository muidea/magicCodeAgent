package parser

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"testing"
)

func TestParse(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset

	/*
		src := `package foo

	import (
		"fmt"
		"time"
	)

	type Model struct {
		Name   string  `json:`
		Desc   float
		Icon   int
		Action []string
	}

	func bar() {
		fmt.Println(time.Now())
	}`


	*/
	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "./module.go", nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		fmt.Println(s.Path.Value)
	}

	for _, s := range f.Scope.Objects {
		fmt.Println(s.Name)
	}

	fmt.Println("abab")
}

func TestParseObject(t *testing.T) {
	packagePtr, packageErr := Parse("./module.go")
	if packageErr != nil {
		t.Errorf("parse failed, err:%s", packageErr.Error())
		return
	}

	byteContent, byteErr := json.Marshal(packagePtr)
	if byteErr != nil {
		t.Errorf("marshal failed, err:%s", byteErr.Error())
		return
	}
	fmt.Print(string(byteContent))
}
