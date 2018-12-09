package common

import "flag"

// MasterOptions contains the command-line options passed to the master process
type MasterOptions struct {
	Program             string // The program to run on the worker nodes
	ConfigPath          string // The path to the config file
	TargetNumExecutions int    // The number of executions to target
}

// NewMasterOptions parses the command-line flags for starting a new master process and stores them in an
// instance of MasterOptions
func NewMasterOptions(program string) *MasterOptions {
	options := new(MasterOptions)
	options.Program = program
	flag.StringVar(&options.ConfigPath, "config", "GoPipeline.config.yaml", "The path to a config file")
	flag.IntVar(&options.TargetNumExecutions, "num", 100, "The number of executions to target")
	flag.Parse()
	return options
}

// WorkerOptions contains the command-line options passed to the worker process
type WorkerOptions struct {
	MasterAddress string // The internet address of the master node
	Position      int    // The position of the worker process within the pipeline stages
	StageID       string // The ID of the stage being run by this worker
}

// NewWorkerOptions parses the command-line flags for starting a new worker process and stores them in an
// instance of WorkerOptions
func NewWorkerOptions() *WorkerOptions {
	options := new(WorkerOptions)
	flag.StringVar(&options.MasterAddress, "address", "127.0.0.1",
		"The internet address of the node running the master process")
	flag.StringVar(&options.StageID, "id", "", "The ID of the stage to be executed")
	flag.IntVar(&options.Position, "position", 0, "The position of the worker process within the pipeline stages")
	flag.Parse()
	return options
}
