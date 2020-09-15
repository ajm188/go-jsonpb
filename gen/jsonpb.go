package gen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// JSON2Marshaler returns a function declaration of an instance method on the
// type named by `id`. The instance method will define a function like:
//
// 		func (m *T) MarshalJSON() ([]byte, error) {
// 			return jsonpb.Marshal(m)
//		}
// where `jsonpb` is the package defined in github.com/ajm188/go-jsonpb
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
									Name: "jsonpb",
								},
								Sel: &ast.Ident{
									Name: "Marshal",
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

// Generate parses a go file in `src`, adds an import of
// github.com/ajm188/go-jsonpb (our jsonpb wrapper) adds a `json.Marshaler`
// implementation to every protobuf struct in the file, and writes the transformed
// AST to `dest`.
func Generate(src string, dest io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return err
	}

	if ok := astutil.AddImport(fset, f, "github.com/ajm188/go-jsonpb"); !ok {
		fmt.Printf("Failed to add json2 import.\nPlease add the following line:\n\t")
		fmt.Println(`import "github.com/ajm188/go-jsonpb"`)
	}

	n := astutil.Apply(f, PreApply, PostApply)
	if err := format.Node(dest, fset, n); err != nil {
		return fmt.Errorf("error dumping ast: %s", err)
	}

	return nil
}
