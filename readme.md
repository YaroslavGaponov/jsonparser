JSON stream parser
==========


# Example

## Code

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/YaroslavGaponov/jsonparser"
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
```

## Result

```sh
Array start
Object start
Key = _id
Value  string 66cf3d8586e785a0362da88f
Key = index
Value  int 0
Key = guid
Value  string 00fa251f-35a0-48b8-8104-30ae5ab4b434
Key = isActive
Value  bool true
Key = balance
Value  string $3522.79
Key = picture
Value  string http://placehold.it/32x32
Key = age
Value  int 20
Key = eyeColor
Value  string green
Key = name
Value  string JeannieDalton
Key = gender
```