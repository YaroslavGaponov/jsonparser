package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"../../jsonparser"
)

func main() {

	file, err := os.Open("example.json")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	parser := jsonparser.New(reader)

	parser.OnObjectStart(func() {
		fmt.Println("Object start")
	})

	parser.OnObjectEnd(func() {
		fmt.Println("Object end")
	})

	parser.OnArrayStart(func() {
		fmt.Println("Array start")
	})

	parser.OnArrayEnd(func() {
		fmt.Println("Array end")
	})

	parser.OnKey(func(name string) {
		fmt.Println("Key =", name)
	})

	parser.OnValue(func(value interface{}) {
		fmt.Println("Value ", reflect.TypeOf(value) , value)
	})

	if err := parser.Run(); err != nil {
		fmt.Println(err)
	}

}
