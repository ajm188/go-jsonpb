package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []string{
		"single_rpc",
		"mixed_structs_and_protos",
		"already_defined",
	}

	skip := map[string]string{
		"mixed_structs_and_protos": "",
		"already_defined":          "",
	}

	for _, testname := range tests {
		name := testname
		t.Run(name, func(t *testing.T) {
			reason, ok := skip[name]
			if ok {
				t.Skip(reason)
			}

			f, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s-in.pb.go.txt", name))
			if err != nil {
				t.Errorf("cannot read source file for %s: %s\n", name, err)
				return
			}

			expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s-out.pb.go.txt", name))
			if err != nil {
				t.Errorf("cannot read dest file for %s: %s", name, err)
			}

			b := &bytes.Buffer{}

			if err := Generate(string(f), b); err != nil {
				t.Errorf("Generate() error: %s\n", err)
				return
			}

			if reflect.DeepEqual(b.Bytes(), expected) {
				t.Errorf("Generate() mismatch:\nhave\n%s\n\nwant\n%s\n", b.Bytes(), expected)
			}
		})
	}
}
