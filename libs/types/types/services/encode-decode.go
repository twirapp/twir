package services

import (
	"bytes"
	"encoding/gob"
)

func Encode[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	return buf.Bytes(), err
}

// nolint:ireturn
func Decode[T any](msg []byte) (T, error) {
	var data T
	err := gob.NewDecoder(bytes.NewReader(msg)).Decode(&data)
	return data, err
}
