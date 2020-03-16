package main

import (
	"bytes"
	"fmt"
	"go/ast"
	_ "go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ExtractFunc extracts information from the provided TypeSpec and returns true if the type should be
// removed from the destination file.
type ExtractFunc func(*ast.TypeSpec) bool

// OptionalFunc returns true if the provided local name is a type that has protobuf.nullable=true
// and should have its marshal functions adjusted to remove the 'Items' accessor.
type OptionalFunc func(name string) bool

// Visitor ..
type Visitor int

// Visit returns a ast visitor
func (v Visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

func main() {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "/Users/xsky/go/src/k8s.io/gengo/examples/error-code/inputs/rbd_error.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)

	ast.Inspect(f, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			fmt.Printf("return statement found on line %v:\n", fset.Position(ret.Pos()))
			printer.Fprint(os.Stdout, fset, ret)
			fmt.Printf("\n")
			return true
		}
		if ret, ok := n.(*ast.GenDecl); ok {
			if ret.Doc != nil && len(ret.Doc.List) > 0 {
				comments := ret.Doc.List
				for _, comment := range comments {
					fmt.Printf("***comment = %v, ", comment)
				}
				fmt.Printf("\n")
				// extractTag("+", comments)
			}
		}
		return true
	})

	// var v Visitor
	// ast.Walk(v, f)

	err = rewriteFile("/Users/xsky/go/src/xsky-demon/external/rbd_error.go", []byte{}, func(*token.FileSet, *ast.File) error { return nil })
	if err != nil {
		panic(err)
	}
}

// RewriteGeneratedGogoProtobufFile generates protobuf file
func RewriteGeneratedGogoProtobufFile(name string, extractFn ExtractFunc, optionalFn OptionalFunc, header []byte) error {

	return rewriteFile(name, header, func(fset *token.FileSet, file *ast.File) error {
		cmap := ast.NewCommentMap(fset, file, file.Comments)

		// transform methods that point to optional maps or slices
		// for _, d := range file.Decls {
		// rewriteOptionalMethods(d, optionalFn)
		// }

		// remove types that are already declared
		decls := []ast.Decl{}
		// for _, d := range file.Decls {
		// if dropExistingTypeDeclarations(d, extractFn) {
		// 	continue
		// }
		// if dropEmptyImportDeclarations(d) {
		// 	continue
		// }
		// decls = append(decls, d)
		// }
		file.Decls = decls
		// remove unmapped comments
		file.Comments = cmap.Filter(file).Comments()
		return nil
	})
}

func rewriteFile(name string, header []byte, rewriteFn func(*token.FileSet, *ast.File) error) error {

	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fset, name, src, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	if err := rewriteFn(fset, file); err != nil {
		return err
	}

	b := &bytes.Buffer{}
	b.Write(header)
	if err := printer.Fprint(b, fset, file); err != nil {
		return err
	}

	body, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(name+".ast", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(body); err != nil {
		return err
	}
	return f.Close()
}

// extractTag gets the comment-tags for the key.  If the tag did not exist, it
// returns the empty string.
func extractTag(key string, lines []string) string {

	val, present := ExtractCommentTags("+", lines)[key]
	if !present || len(val) < 1 {
		return ""
	}
	return val[0]
}

// ExtractCommentTags parses comments for lines of the form:
//
//   'marker' + "key=value".
//
// Values are optional; "" is the default.  A tag can be specified more than
// one time and all values are returned.  If the resulting map has an entry for
// a key, the value (a slice) is guaranteed to have at least 1 element.
//
// Example: if you pass "+" for 'marker', and the following lines are in
// the comments:
//   +foo=value1
//   +bar
//   +foo=value2
//   +baz="qux"
// Then this function will return:
//   map[string][]string{"foo":{"value1, "value2"}, "bar": {""}, "baz": {"qux"}}
func ExtractCommentTags(marker string, lines []string) map[string][]string {
	out := map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, marker) {
			continue
		}
		// TODO: we could support multiple values per key if we split on spaces
		kv := strings.SplitN(line[len(marker):], "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = append(out[kv[0]], kv[1])
		} else if len(kv) == 1 {
			out[kv[0]] = append(out[kv[0]], "")
		}
	}
	return out
}

// ExtractSingleBoolCommentTag parses comments for lines of the form:
//
//   'marker' + "key=value1"
//
// If the tag is not found, the default value is returned.  Values are asserted
// to be boolean ("true" or "false"), and any other value will cause an error
// to be returned.  If the key has multiple values, the first one will be used.
func ExtractSingleBoolCommentTag(marker string, key string, defaultVal bool, lines []string) (bool, error) {
	values := ExtractCommentTags(marker, lines)[key]
	if values == nil {
		return defaultVal, nil
	}
	if values[0] == "true" {
		return true, nil
	}
	if values[0] == "false" {
		return false, nil
	}
	return false, fmt.Errorf("tag value for %q is not boolean: %q", key, values[0])
}

func gopathDir(pkg string) (string, error) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		absPath, err := filepath.Abs(path.Join(gopath, "src", pkg))
		if err != nil {
			return "", err
		}
		if dir, err := os.Stat(absPath); err == nil && dir.IsDir() {
			return absPath, nil
		}
	}
	return "", fmt.Errorf("%s not in $GOPATH", pkg)
}

func (i *customImporter) fsPkg(pkg string) (*types.Package, error) {
	dir, err := gopathDir(pkg)
	if err != nil {
		return importOrErr(i.base, pkg, err)
	}

	dirFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return importOrErr(i.base, pkg, err)
	}

	fset := token.NewFileSet()
	var files []*ast.File
	for _, fileInfo := range dirFiles {
		if fileInfo.IsDir() {
			continue
		}
		n := fileInfo.Name()
		if path.Ext(fileInfo.Name()) != ".go" {
			continue
		}
		if i.skipTestFiles && strings.Contains(fileInfo.Name(), "_test.go") {
			continue
		}
		file := path.Join(dir, n)
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		f, err := parser.ParseFile(fset, file, src, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	conf := types.Config{
		Importer: i,
	}
	p, err := conf.Check(pkg, fset, files, nil)

	if err != nil {
		return importOrErr(i.base, pkg, err)
	}
	return p, nil
}

func importOrErr(base types.Importer, pkg string, err error) (*types.Package, error) {
	p, impErr := base.Import(pkg)
	if impErr != nil {
		return nil, err
	}
	return p, nil
}

type customImporter struct {
	imported      map[string]*types.Package
	base          types.Importer
	skipTestFiles bool
}

func (i *customImporter) Import(path string) (*types.Package, error) {
	var err error
	if path == "" || path[0] == '.' {
		path, err = filepath.Abs(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		path = StripGopath(path)
	}
	if pkg, ok := i.imported[path]; ok {
		return pkg, nil
	}
	pkg, err := i.fsPkg(path)
	if err != nil {
		return nil, err
	}
	i.imported[path] = pkg
	return pkg, nil
}

// StripGopath teks the directory to a package and remove the gopath to get the
// cannonical package name
func StripGopath(p string) string {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		p = strings.Replace(p, path.Join(gopath, "src")+"/", "", 1)
	}
	return p
}
