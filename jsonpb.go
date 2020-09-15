package jsonpb

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func Marshal(pb proto.Message) ([]byte, error) {
	buf := bytes.Buffer{}
	m := jsonpb.Marshaler{}

	if err := m.Marshal(buf, pb); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
