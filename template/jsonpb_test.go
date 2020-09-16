package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []string{
		"single_rpc",
		"nested_structs",
	}

	skip := map[string]string{}

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

			expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s-out.go.txt", name))
			if err != nil {
				t.Errorf("cannot read dest file for %s: %s", name, err)
			}

			b := &bytes.Buffer{}

			if err := Generate(string(f), b); err != nil {
				t.Errorf("Generate() error: %s\n", err)
				return
			}

			if string(b.Bytes()) != string(expected) {
				t.Errorf("Generate() mismatch:\nhave\n----\n%s\nwant\n%s\n", b.Bytes(), expected)
			}
		})
	}
}
