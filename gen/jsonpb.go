package gen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
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
									Name: "protojson",
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

// PreApply receives a cursor, which maintains its location in the parse tree.
// This is our pre-apply function, which means that if this function returns `true`,
// then PreApply will be called for that Node's children, and the post-apply
// function will be called for that Node.
// See astutil.Apply (https://pkg.go.dev/golang.org/x/tools/go/ast/astutil#Apply) for more.
func PreApply(cur *astutil.Cursor) bool {
	if cur.Node() == nil {
		return false
	}

	switch node := cur.Node().(type) {
	case *ast.File:
		// File will always be the root node of a parse tree for our purposes,
		// so if we return false here, then no child nodes will be processed at all,
		// instead of proceeding downward to find the struct declarations.
		return true
	case *ast.GenDecl:
		// GenDecl is one of (import var, const, type). We care only about the `type`
		// variant.
		//
		// Note that functions have their own Decl implementation, namely FuncDecl.
		return node.Tok == token.TYPE
	default:
	}

	log.Printf("(pre) found %T, not traversing it\n", cur.Node())

	return false
}

// NewPostApplier returns an astutil.ApplyFunc which, when it encounters a struct
// declaration where that struct has any struct tag containing the text "protobuf:"
// calls the passed in protofunc. This allows you to do something like, track the
// names of all proto structs.
func NewPostApplier(protofunc func(*astutil.Cursor, *ast.TypeSpec)) astutil.ApplyFunc {
	return func(cur *astutil.Cursor) bool {
		if cur.Node() == nil {
			return true
		}

		switch node := cur.Node().(type) {
		case *ast.File:
			// return value doesn't strictly matter here, since this will be the last node to
			// be post-applied.
			return false
		case *ast.GenDecl:
			// GenDecl is one of (import, var, const, type).
			//
			// GenDecls can have multiple specs (definitions) within a single declaration. Think syntax like:
			//		type (
			// 			A struct{}
			//			B struct{}
			//			C interface{}
			//		)
			// So, we need to iterate, but on each iteration assert that the Spec is indeed a TypeSpec.
			// Even though it is invalid to write Go code that declares a var and a type in the same Decl,
			// it is possible to represent that in the type system with AST nodes. Better safe than sorry.
			for _, spec := range node.Specs {
				log.Printf("decl spec: %+v (%T)\n", spec, spec)

				switch ts := spec.(type) {
				case *ast.TypeSpec:
					// TypeSpec.Type can be a lot of things (see https://golang.org/pkg/go/ast/#TypeSpec),
					// but in our case we're only interested in structs.
					switch s := ts.Type.(type) {
					case *ast.StructType:
						isProto := false

						if s.Fields.List == nil {
							return true
						}

						// Naive, but good enough, check if this is a proto. If any Field on the struct
						// has a tag that looks like it came from protoc-gen-go, by which we mean containing
						// the string "protobuf:", then it's probably a protobuf struct.
						for _, field := range s.Fields.List {
							if field.Tag == nil {
								continue
							}

							if strings.Contains(field.Tag.Value, "protobuf:") {
								isProto = true
								break
							}
						}

						// Allow callers to specify a callback, if they want to do any custom action on
						// finding a protobuf struct.
						if isProto {
							protofunc(cur, ts)
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
			log.Printf("yikes, got something that's not a GenDecl %T\n", cur.Node())
		}

		return true
	}
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

	if ok := astutil.AddImport(fset, f, "google.golang.org/protobuf/encoding/protojson"); !ok {
		// These messages have to go regardless of log level.
		fmt.Fprintf(os.Stderr, "Failed to add json2 import.\nPlease add the following line:\n\t")
		fmt.Fprintln(os.Stderr, `import "google.golang.org/protobuf/encoding/protojson"`)
	}

	n := astutil.Apply(f, PreApply, NewPostApplier(func(cur *astutil.Cursor, spec *ast.TypeSpec) {
		cur.InsertAfter(JSON2Marshaler(spec.Name))
	}))
	if err := format.Node(dest, fset, n); err != nil {
		return fmt.Errorf("error dumping ast: %s", err)
	}

	return nil
}
