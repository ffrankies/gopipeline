// Package gopipeline contains the logic for pipelining a Golang application through a list of compute nodes.
package gopipeline

import (
	"errors"
	"os"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/master"
	"github.com/ffrankies/gopipeline/types"
	"github.com/ffrankies/gopipeline/worker"
)

// getProcessType obtains the process type from the command line arguments.
func getProcessType() (processType string, err error) {
	invalidArgError := errors.New("Must pass in either \"master\" or \"worker\" as the last command-line argument")
	if len(os.Args) < 2 {
		err = invalidArgError
		return
	}
	processType = os.Args[len(os.Args)-1]
	if processType != "master" && processType != "worker" {
		err = invalidArgError
		return
	}
	return
}

// Run runs either the master or the worker stage on a single node. Also parses the command-line arguments needed for
// the worker and/or master
func Run(functionList []types.AnyFunc, registerType interface{}) {
	program := os.Args[0]
	processType, err := getProcessType()
	if err != nil {
		panic(err)
	}
	if processType == "master" {
		options := common.NewMasterOptions(program)
		master.Run(options, functionList)
		return
	}
	if processType == "worker" {
		options := common.NewWorkerOptions()
		worker.Run(options, functionList, registerType)
		return
	}
}
