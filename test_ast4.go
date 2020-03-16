package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

type visitor func(n ast.Node) ast.Visitor

func (v visitor) Visit(n ast.Node) ast.Visitor {
	return v(n)
}

// errType just defines error in terms of go/types
var errType *types.Interface

func init() {
	errType = types.NewInterfaceType([]*types.Func{
		types.NewFunc(0, nil, "Error",
			types.NewSignature(
				nil,              // Receiver
				types.NewTuple(), // Params
				types.NewTuple(types.NewParam(0, nil, "", types.Typ[types.String])), // Result
				false,
			),
		),
	}, nil)

	errType.Complete()
}

func main() {
	// parse the directory of go code
	fset := token.NewFileSet()
	src, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		panic(err)
	}

	// type check the code
	conf := types.Config{Importer: importer.Default()}
	fs := []*ast.File{}
	for _, tree := range src {
		for _, f := range tree.Files {
			fs = append(fs, f)
		}
	}
	fmt.Println(fs)
	i := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}}
	if _, err := conf.Check("cmd/hello", fset, fs, i); err != nil {
		log.Fatal(err) // type error
	}
	fmt.Println("*****")

	// walk the ast, priting any expression that implements to error
	var v visitor

	v = func(n ast.Node) ast.Visitor {
		e, ok := n.(ast.Expr)
		if !ok {
			return v
		}

		t := i.Types[e].Type
		if implements(t, errType) {
			fmt.Printf("%s, %T %+v %v\n", fset.Position(n.Pos()), n, n, t)
		}

		return v
	}

	ast.Walk(v, fs[0])
}

// implements returns false if t doesn't implement i
//
// The standard types.Implements panics on false, making it inconvenient for
// simple checking.
func implements(t types.Type, i *types.Interface) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	ok = types.Implements(t, i)

	return
}
