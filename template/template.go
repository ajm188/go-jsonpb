package template

const jsonpbTemplate = `package {{ .Name }}

import "github.com/ajm188/go-jsonpb"

{{ range .ProtoTypes -}}
func (m *{{ .Type }}) MarshalJSON() ([]byte, error) {
	return jsonpb.Marshal(m)
}

{{- end }}
`
