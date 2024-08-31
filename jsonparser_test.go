package jsonparser_test

import (
	"strings"
	"testing"

	"./"
)

func TestSimpleValueIsString(t *testing.T) {
	o := `{
		"hello": "welcome"
	}`

	parser := jsonparser.New(strings.NewReader(o))

	result := []string{}

	parser.OnObjectStart(func() {
		result = append(result, "start")
	})

	parser.OnObjectEnd(func() {
		result = append(result, "end")
	})

	parser.OnKey(func(key string) {
		result = append(result, key)
	})

	parser.OnValue(func(value interface{}) {
		str, ok := value.(string)
		if !ok {
			t.Errorf("value is not string")
		}
		result = append(result, str)
	})

	if err := parser.Run(); err != nil {
		t.Fatal(err)
	}

	if len(result) != 4 {
		t.Error("not all handlers are called")
	}
	if result[0] != "start" {
		t.Error("start object is not found")
	}
	if result[1] != "hello" {
		t.Error("key is not found")
	}
	if result[2] != "welcome" {
		t.Error("value is not found")
	}
	if result[3] != "end" {
		t.Error("end object is not found")
	}

}
