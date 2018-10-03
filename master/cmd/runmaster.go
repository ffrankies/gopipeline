package main

import (
	"github.com/ffrankies/gopipeline"
	"github.com/ffrankies/gopipeline/master"
)

func main() {
	master.Run("example", "GoPipeline.config", make([]gopipeline.AnyFunc, 0))
}
