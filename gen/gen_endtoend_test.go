// +build itest

package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ensureProtoc() error {
	_, err := exec.LookPath("protoc")
	return err
}

func compilePlugin() error {
	cmd := exec.Command("go", "build", "-o", "protoc-gen-go-json", "../cmd/protoc-gen-go-json/")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func TestEndToEnd(t *testing.T) {
	// There is currently no way to _not_ have protoc/protogen write to stdout
	// (ref: https://github.com/protocolbuffers/protobuf-go/blob/db5c900f0ce544131509b33f6d68ec651e3ca91c/compiler/protogen/protogen.go#L89)
	// so, these are going to be "integration tests" of sorts that (a) require protoc to be
	// in the path and (b) require a compiled version of protoc-gen-go-json.

	require.NoError(t, ensureProtoc())
	require.NoError(t, compilePlugin())

	// Remove the locally-built binary
	defer os.Remove("./protoc-gen-go-json")

	envPATH := os.Getenv("PATH")

	cmd := exec.Command("protoc", "--go-json_out=.", "-Itestdata", "testdata/test.proto")
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=.:%s", envPATH))

	require.NoError(t, cmd.Run())
	// Cleanup the generated code
	defer os.Remove("test_json.pb.go")

	result, err := ioutil.ReadFile("test_json.pb.go")
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/test_json.pb.go.out")
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(result))
}
