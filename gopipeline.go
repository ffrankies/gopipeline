// Package gopipeline contains the logic for pipelining a Golang application through a list of compute nodes.
package gopipeline

import (
	"flag"
	"fmt"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/master"
)

// Run runs either the master or the worker stage on a single node. Also parses the command-line arguments needed for
// the worker and/or master
func Run(functionList []common.AnyFunc) {
	var processType string
	flag.StringVar(&processType, "type", "master", "Determines the process type. Can be either of [master, worker]")
	if processType == "master" {
		master.Run("hostname", "GoPipeline.config.yaml", functionList)
		return
	}
	if processType == "worker" {
		fmt.Println("Running a worker process! Woohoo!")
	}
}
