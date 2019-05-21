package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"

	"github.com/pkg/errors"
)

func EncodeObject(object interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(object)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encode object: '%+v'", object)
	}

	return buffer.Bytes(), nil
}

func DecodeObject(data []byte, object interface{}) error {
	buffer := bytes.NewReader(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(object)
	if err != nil {
		return errors.Wrapf(err, "failed to decode data: '%# x'", hex.EncodeToString(data))
	}

	return nil
}
