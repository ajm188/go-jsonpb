package jsonpb

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Marshal uses protobuf/jsonpb to marshal a proto message into JSON.
func Marshal(pb proto.Message) ([]byte, error) {
	return protojson.Marshal(pb)
}
