package template

const jsonpbTemplate = `package {{ .Name }}

import "google.golang.org/protobuf/encoding/protojson"
{{ range .ProtoTypes }}
func (m *{{ .Type }}) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(m)
}
{{ end }}`
