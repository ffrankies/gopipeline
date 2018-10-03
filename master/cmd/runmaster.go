package main

import (
	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/master"
)

func main() {
	master.Run("README.md", make([]gopipeline.AnyFunc, 0))
}
