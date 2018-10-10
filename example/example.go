// This is an example package that uses the gopipeline library
package main

import (
	"fmt"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/types"
)

func hello(args ...interface{}) interface{} {
	fmt.Println("Hello World")
	return "wut"
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	for i := 0; i < 5; i++ {
		functionList = append(functionList, hello)
	}
	gopipeline.Run(functionList)
}
