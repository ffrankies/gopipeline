// Package gopipeline contains the logic for pipelining a Golang application through a list of compute nodes.
package gopipeline

import (
	"fmt"
	"os"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/master"
	"github.com/ffrankies/gopipeline/types"
)

// Run runs either the master or the worker stage on a single node. Also parses the command-line arguments needed for
// the worker and/or master
func Run(functionList []types.AnyFunc) {

	processType := os.Args[1]
	// flag.StringVar(&processType, "type", "master", "Determines the process type. Can be either of [master, worker]")
	fmt.Println("Got a command line argument...", processType)
	return
	if processType == "master" {
		options := common.NewMasterOptions()
		master.Run(options, functionList)
		return
	}
	if processType == "worker" {
		options := common.NewWorkerOptions()
		fmt.Println("Worker options =", options)
		fmt.Println("Running a worker process! Woohoo!")
	}
}
