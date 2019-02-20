package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	A int
	B string
	C float32
	D bool
}

type TestObject2 struct {
	A int `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	C float32 `json:"c,omitempty"`
	D bool
}

type SerializeTestInputOutput struct {
	expected string
	obj interface{}
}

func TestSerializeOfJSONSerializer(t *testing.T) {
	tc := []SerializeTestInputOutput{
		0: SerializeTestInputOutput{
			expected: `{"A":22,"B":"bbbb","C":5.25,"D":false}`,
			obj: TestObject{
				A: 22,
				B: "bbbb",
				C: 5.25,
				D: false,
			},
		},
		1: SerializeTestInputOutput{
			expected: `{"A":0,"B":"","C":0,"D":false}`,
			obj: TestObject{},
		},
		2: SerializeTestInputOutput{
			expected: `{"D":false}`,
			obj: TestObject2{},
		},
		3: SerializeTestInputOutput{
			expected: `{"a":55,"D":true}`,
			obj: TestObject2{
				A: 55,
				D: true,
			},
		},
	}
	s := JSONSerializer{}
	for _, item := range tc {
		bytes, _ := s.Serialize(item.obj)
		assert.Equal(t, item.expected, string(bytes))
	}
}

func TestDeserializeOfJSONSerializer(t *testing.T) {
	s := JSONSerializer{}
	var actualObject TestObject
	s.Deserialize(strings.NewReader(`{"A":22,"B":"bbbb","C":5.25,"D":false}`), &actualObject)
	obj := TestObject{
		A: 22,
		B: "bbbb",
		C: 5.25,
		D: false,
	}
	assert.Equal(t, obj, actualObject)
}