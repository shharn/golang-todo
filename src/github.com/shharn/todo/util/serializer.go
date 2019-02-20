package util

import (
	"io"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize(io.Reader, interface{}) error
}

type JSONSerializer struct {}

func (js JSONSerializer) Serialize(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, nil
	}

	if bytes, err := json.Marshal(obj); err == nil {
		return bytes, nil
	} else {
		return nil, errors.WithStack(err)
	}
}

func (js JSONSerializer) Deserialize(reader io.Reader, target interface{}) error {
	if reader == nil {
		return nil
	}

	if !js.isPointerType(target) {
		return errors.New("Object you want to serialize should be pointer type")
	}

	if err := json.NewDecoder(reader).Decode(target); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (js JSONSerializer) isPointerType(obj interface{}) bool {
	tp := reflect.TypeOf(obj)
	return tp.Kind() == reflect.Ptr
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}