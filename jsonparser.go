package jsonparser

import (
	"fmt"
	"io"
)

type JSONType uint8

const (
	Object JSONType = iota
	Array
	Null
	String
	Number
	Boolean
)

type KeyValue bool

const (
	Key   KeyValue = false
	Value KeyValue = true
)

type JSONParser struct {
	stream io.RuneReader

	path     []JSONType
	buffer   []rune
	isString bool
	keyvalue KeyValue

	onStartObject func()
	onEndObject   func()

	onStartArray func()
	onEndArray   func()

	onKey   func(name string)
	onValue func(value interface{})
}

func New(stream io.RuneReader) *JSONParser {
	return &JSONParser{
		stream:   stream,
		path:     make([]JSONType, 0, 5),
		keyvalue: Key,
		isString: false,
		buffer:   make([]rune, 0, 512),
	}
}

func (parser *JSONParser) OnObjectStart(handler func()) {
	parser.onStartObject = handler
}

func (parser *JSONParser) OnObjectEnd(handler func()) {
	parser.onEndObject = handler
}

func (parser *JSONParser) OnArrayStart(handler func()) {
	parser.onStartArray = handler
}

func (parser *JSONParser) OnArrayEnd(handler func()) {
	parser.onEndArray = handler
}

func (parser *JSONParser) OnKey(handler func(name string)) {
	parser.onKey = handler
}

func (parser *JSONParser) OnValue(handler func(value interface{})) {
	parser.onValue = handler
}

func (parser *JSONParser) Run() error {

	for {
		r, _, err := parser.stream.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading rune: %s", err)
		}

		switch string(r) {
		case "{":
			parser.path = append(parser.path, Object)
			parser.onStartObject()
			parser.keyvalue = Key

		case "}":
			parser.path = parser.path[0 : len(parser.path)-1]
			parser.onEndObject()
			parser.keyvalue = Key

		case "[":
			parser.path = append(parser.path, Array)
			parser.onStartArray()
			parser.keyvalue = Value

		case "]":
			parser.path = parser.path[0 : len(parser.path)-1]
			parser.onEndArray()
			parser.keyvalue = Key

		case ":":
			if !parser.isString {
				parser.flush()
				parser.keyvalue = Value
			} else {
				parser.buffer = append(parser.buffer, r)
			}

		case "\"":
			parser.flush()
			parser.isString = !parser.isString

		case " ":
		case "\t":
		case "\n":
		case ",":
			if !parser.isString {
				parser.flush()
			}

		default:
			parser.buffer = append(parser.buffer, r)

		}
	}
	return nil
}

func (parser *JSONParser) flush() {
	if len(parser.buffer) > 0 {
		if parser.keyvalue == Key {
			parser.onKey(string(parser.buffer))
			parser.keyvalue = Value
		} else {
			parser.onValue(string(parser.buffer))
			if parser.path[len(parser.path)-1] == Object {
				parser.keyvalue = Key
			} else {
				parser.keyvalue = Value
			}
		}
		parser.buffer = []rune{}
	}
}
