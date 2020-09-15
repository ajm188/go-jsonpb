package jsonpb

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// Marshal uses protobuf/jsonpb to marshal a proto message into JSON.
func Marshal(pb proto.Message) ([]byte, error) {
	buf := bytes.Buffer{}
	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, pb); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
