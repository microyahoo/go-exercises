package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

func main() {
	src := `package main
        type Example struct {
    Foo string` + " `json:\"foo\"` }"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			fmt.Println(field.Names)
			fmt.Printf("\nField: %s\n", field.Names[0].Name)
			fmt.Printf("Tag:   %s\n", field.Tag.Value)
		}
		return false
	})
}

type config struct {
	// first section - input & output
	file     string
	modified io.Reader
	output   string
	write    bool

	// second section - struct selection
	offset     int
	structName string
	line       string
	start, end int

	// third section - struct modification
	remove     []string
	add        []string
	override   bool
	transform  string
	sort       bool
	clear      bool
	addOpts    []string
	removeOpts []string
	clearOpt   bool
}
