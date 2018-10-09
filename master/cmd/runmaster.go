package main

import (
	"github.com/ffrankies/gopipeline/master"
	"github.com/ffrankies/gopipeline/types"
)

func exampleFunc(args ...interface{}) interface{} {
	a := args[0].(int)
	a += 2
	return a
}

func main() {
	functionList := make([]types.AnyFunc, 0)
	for i := 0; i < 10; i++ {
		functionList = append(functionList, exampleFunc)
	}
	master.Run("example", "GoPipeline.config.yaml", functionList)
}
