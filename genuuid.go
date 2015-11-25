package main

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	fmt.Printf("filenames := {")
	for i := 0; i < 100000; i++ {
		uuid := uuid.NewV4().String()
		fmt.Printf("%s, ", uuid)
	}
	fmt.Printf("}")
}
