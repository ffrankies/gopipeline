package main

import (
	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/master"
)

func exampleFunc(args ...interface{}) interface{} {
	a := args[0].(int)
	a += 2
	return a
}

func main() {
	functionList := make([]gopipeline.AnyFunc, 0)
	for i := 0; i < 10; i++ {
		functionList = append(functionList, exampleFunc)
	}
	master.Run("example", "GoPipeline.config2.yaml", functionList)
}
