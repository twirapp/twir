package protoutils

import (
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func CreateJsonWithProto(msg proto.Message, additionalFields map[string]any) (
	json.RawMessage,
	error,
) {
	var result map[string]any

	protoBytes, err := protojson.MarshalOptions{
		UseEnumNumbers:    true,
		EmitDefaultValues: true,
		EmitUnpopulated:   true,
	}.Marshal(msg)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(protoBytes, &result); err != nil {
		return nil, err
	}

	for k, v := range additionalFields {
		result[k] = v
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return b, err
}
