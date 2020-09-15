package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func getSrc(filename string) string {
	f, err := os.Open(filename)
	must(err)

	defer f.Close()

	src, err := ioutil.ReadAll(f)
	must(err)

	return string(src)
}

func getDest(filename string) *os.File {
	f, err := os.Create(filename)

	must(err)
	must(f.Chmod(0644))

	return f
}

func JSON2Marshaler(id *ast.Ident) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "m",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: id.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "MarshalJSON",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: "byte",
							},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "json2",
								},
								Sel: &ast.Ident{
									Name: "MarshalPB",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{Name: "m"},
							},
						},
					},
				},
			},
		},
	}
}

func PreApply(cur *astutil.Cursor) bool {
	if cur.Node() == nil {
		return false
	}

	switch node := cur.Node().(type) {
	case *ast.File:
		return true
	case *ast.GenDecl:
		return node.Tok == token.TYPE
	default:
	}

	fmt.Printf("(pre) found %T, not traversing it\n", cur.Node())

	return false
}

func PostApply(cur *astutil.Cursor) bool {
	fmt.Printf("(post-apply) %T\n", cur.Node())

	if cur.Node() == nil {
		return true
	}

	switch node := cur.Node().(type) {
	case *ast.File:
		return false
	case *ast.GenDecl:
		for _, spec := range node.Specs {
			fmt.Printf("decl spec: %+v (%T)\n", spec, spec)

			switch ts := spec.(type) {
			case *ast.TypeSpec:
				switch s := ts.Type.(type) {
				case *ast.StructType:
					isProto := false

					if s.Fields.List == nil {
						return true
					}

					for _, field := range s.Fields.List {
						if field.Tag == nil {
							continue
						}

						if strings.Contains(field.Tag.Value, "protobuf:") {
							isProto = true
							break
						}
					}

					if isProto {
						fun := JSON2Marshaler(ts.Name)
						cur.InsertAfter(fun)
					}
				default:
					return true
				}
			default:
				return true
			}
		}
		return true
	default:
		fmt.Printf("yikes, got something that's not a GenDecl %T\n", cur.Node())
	}

	return true
}

func main() {
	srcname := os.Args[1]
	destname := os.Args[2]
	src := getSrc(srcname)
	dest := getDest(destname)

	defer dest.Close()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	must(err)

	if ok := astutil.AddImport(fset, f, "vitess.io/vitess/go/json2"); !ok {
		fmt.Printf("Failed to add json2 import to %s.\nPlease add the following line:\n\t", srcname)
		fmt.Println(`import "vitess.io/vitess/go/json2"`)
	}

	n := astutil.Apply(f, PreApply, PostApply)
	if err := format.Node(dest, fset, n); err != nil {
		fmt.Printf("error dumping ast: %s\n", err)
	}
}
