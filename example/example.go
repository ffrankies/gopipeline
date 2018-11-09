// This is an example package that uses the gopipeline library
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/types"
)

func hello(args ...interface{}) interface{} {
	time.Sleep(2 * time.Second)
	if len(args) > 0 {
		num := args[0].(int)
		fmt.Println("Hello World from Position: " + strconv.Itoa(num))
		num++
		return num
	}
	fmt.Println("Hello World from Position: 0")
	return 1
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	for i := 0; i < 5; i++ {
		functionList = append(functionList, hello)
	}
	gopipeline.Run(functionList)
}
