package buscore

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/goccy/go-json"
)

type QueueEncoder string

const (
	JsonEncoder QueueEncoder = "json"
	GobEncoder  QueueEncoder = "gob"
)

func EncodeGob[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	return buf.Bytes(), err
}

// nolint:ireturn
func DecodeGob[T any](msg []byte) (T, error) {
	var data T
	err := gob.NewDecoder(bytes.NewReader(msg)).Decode(&data)
	return data, err
}

func natsEncode[T any](encoder QueueEncoder, data T) ([]byte, error) {
	switch encoder {
	case JsonEncoder:
		return json.Marshal(data)
	case GobEncoder:
		return EncodeGob(data)
	}

	return nil, fmt.Errorf("unknown encoder")
}

func natsDecode[T any](encoder QueueEncoder, data []byte) (T, error) {
	if len(data) == 0 {
		var emptyRes T
		return emptyRes, nil
	}

	switch encoder {
	case JsonEncoder:
		var res T
		if err := json.Unmarshal(data, &res); err != nil {
			return res, err
		}

		return res, nil
	case GobEncoder:
		return DecodeGob[T](data)
	}

	var emptyRes T
	return emptyRes, fmt.Errorf("unknown decoder")
}
