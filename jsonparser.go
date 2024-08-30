package jsonparser

import (
	"fmt"
	"io"
	"strconv"
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

	delimiters map[rune]interface{}

	path       []JSONType
	buffer     []rune
	flagString bool
	keyvalue   KeyValue

	onStartObject func()
	onEndObject   func()

	onStartArray func()
	onEndArray   func()

	onKey   func(name string)
	onValue func(value interface{})
}

func New(stream io.RuneReader) *JSONParser {

	return &JSONParser{
		stream: stream,

		delimiters: map[rune]interface{}{' ': true, ':': true, '\n': true, '\t': true, ',': true},

		path:       make([]JSONType, 0, 10),
		keyvalue:   Key,
		flagString: false,
		buffer:     make([]rune, 0, 512),

		onStartObject: func() {},
		onEndObject:   func() {},

		onStartArray: func() {},
		onEndArray:   func() {},

		onKey:   func(name string) {},
		onValue: func(value interface{}) {},
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
			if !parser.flagString {
				parser.flush()
				parser.keyvalue = Value
			} else {
				parser.buffer = append(parser.buffer, r)
			}

		case "\"":
			parser.flush()
			parser.flagString = !parser.flagString

		default:
			if parser.flagString {
				parser.buffer = append(parser.buffer, r)
			} else {
				if _, ok := parser.delimiters[r]; ok {
					parser.flush()
				} else {
					parser.buffer = append(parser.buffer, r)
				}
			}
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
			parser.onValue(cast(parser.buffer))
			if parser.path[len(parser.path)-1] == Object {
				parser.keyvalue = Key
			} else {
				parser.keyvalue = Value
			}
		}
		parser.buffer = parser.buffer[:0]
	}
}

func cast(value []rune) interface{} {
	s := string(value)
	switch s {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return Null
	default:
		if num, err := strconv.Atoi(s); err == nil {
			return num
		}
		return s
	}
}
