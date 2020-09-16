package template

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"text/template"

	"github.com/ajm188/go-jsonpb/gen"
	"golang.org/x/tools/go/ast/astutil"
)

type ProtoPackage struct {
	Name       string
	ProtoTypes []ProtoType
}

type ProtoType struct {
	Type string
}

func Generate(src string, dest io.Writer) error {
	fset := token.NewFileSet()
	// We don't care about comments when doing template-generation mode, since
	// the original file will be untouched.
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return err
	}

	tmpl, err := template.New("jsonpb").Parse(jsonpbTemplate)
	if err != nil {
		return err
	}

	pkg := ProtoPackage{}

	preapplier := func(cur *astutil.Cursor) bool {
		if cur.Node() == nil {
			return true
		}

		file, ok := cur.Node().(*ast.File)
		if !ok {
			return true
		}

		pkg.Name = file.Name.Name

		return true
	}

	postapplier := func(cur *astutil.Cursor, spec *ast.TypeSpec) {
		pkg.ProtoTypes = append(pkg.ProtoTypes, ProtoType{spec.Name.Name})
	}

	_ = astutil.Apply(f, preapplier, gen.NewPostApplier(postapplier))

	return tmpl.Execute(dest, pkg)
}
