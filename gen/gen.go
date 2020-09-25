package gen

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// GeneratedFile is the set of methods needed by protoc-gen-go-json
// on *protogen.GeneratedFile. It exists mainly for ease of testing.
type GeneratedFile interface {
	P(...interface{})
}

func GenerateFile(g GeneratedFile, req *pluginpb.CodeGeneratorRequest, file *protogen.File) {
	GenerateHeader(g, req, file)

	g.P("package ", file.GoPackageName)
	g.P()

	g.P(`import "google.golang.org/protobuf/encoding/protojson"`)
	g.P()

	GenerateMarshalers(g, collectMessages(file))
}

func GenerateHeader(g GeneratedFile, req *pluginpb.CodeGeneratorRequest, file *protogen.File) {
	protocVersion := "(unknown)"
	if v := req.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}

	g.P("// Code generated by protoc-gen-go-json. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// \tprotoc-gen-go-json ", "(unknown)")
	g.P("// \tprotoc             ", protocVersion)

	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}

	g.P()
}

func GenerateMarshalers(g GeneratedFile, messages []*protogen.Message) {
	for _, message := range messages {
		g.P("func (x *", message.GoIdent, ") MarshalJSON() ([]byte, error) {")
		g.P("return protojson.Marshal(x)")
		g.P("}")
		g.P()
	}
}

// Implementation lifted nearly wholesale from protoc-gen-go.
// https://github.com/protocolbuffers/protobuf-go/blob/db5c900f0ce544131509b33f6d68ec651e3ca91c/cmd/protoc-gen-go/internal_gengo/init.go#L50-L84
func collectMessages(file *protogen.File) []*protogen.Message {
	allMessages := []*protogen.Message{}

	var walk func([]*protogen.Message, func(*protogen.Message))
	walk = func(messages []*protogen.Message, f func(*protogen.Message)) {
		for _, m := range messages {
			f(m)
			walk(m.Messages, f)
		}
	}

	initMessages := func(messages []*protogen.Message) {
		for _, message := range messages {
			allMessages = append(allMessages, message)
		}
	}

	initMessages(file.Messages)
	walk(file.Messages, func(m *protogen.Message) {
		initMessages(m.Messages)
	})

	return allMessages
}
