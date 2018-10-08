package main

import (
	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/master"
)

func exampleFunc(args ...interface{}) interface{} {
	a := args[0].(int)
	a += 2
	return a
}

func main() {
	functionList := make([]common.AnyFunc, 0)
	for i := 0; i < 10; i++ {
		functionList = append(functionList, exampleFunc)
	}
	master.Run("example", "GoPipeline.config.yaml", functionList)
}
