package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// borrowed from https://github.com/lukehoban/go-outline/blob/master/main.go#L54-L107
	fset := token.NewFileSet()
	parserMode := parser.ParseComments
	var fileAst *ast.File
	var err error
	fileAst, err = parser.ParseFile(fset, "group.go", nil, parserMode)
	if err != nil {
		panic(err)
	}
	spew.Dump(fileAst)
	printer.Fprint(os.Stdout, fset, fileAst)
	for _, d := range fileAst.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			fmt.Println("Func")
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println("Import", spec.Path.Value)
				case *ast.TypeSpec:
					fmt.Println("Type", spec.Name.String())
				case *ast.ValueSpec:
					for _, id := range spec.Names {
						fmt.Printf("id.Obj.Decl = %+v, type(Decl) = %T, Data = %#v, type(Data) = %T,  Kind = %v, Type = %#v, Name=%v\n", id.Obj.Decl, id.Obj.Decl, id.Obj.Data, id.Obj.Data, id.Obj.Kind, id.Obj.Type, id.Obj.Name)
						if d, ok := id.Obj.Decl.(*ast.ValueSpec); ok {
							fmt.Printf("*** valueSpec = %+v\n", d)
						}
						// fmt.Printf("Var %s: %v", id.Name, id.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value)
					}
				default:
					fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			fmt.Printf("Unknown declaration @\n", decl.Pos())
		}
	}
}
