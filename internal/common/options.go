package common

import "flag"

// MasterOptions contains the command-line options passed to the master process
type MasterOptions struct {
	Program    string // The program to run on the worker nodes
	ConfigPath string // The path to the config file
}

// NewMasterOptions parses the command-line flags for starting a new master process and stores them in an
// instance of MasterOptions
func NewMasterOptions() *MasterOptions {
	options := new(MasterOptions)
	flag.StringVar(&options.Program, "program", "example", "The program to run on the worker nodes")
	flag.StringVar(&options.ConfigPath, "config", "GoPipeline.config.yaml", "The path to a config file")
	return options
}

// WorkerOptions contains the command-line options passed to the worker process
type WorkerOptions struct {
	MasterAddress string // The internet address of the master node
	MasterPort    int    // The port number of the master process
	Position      int    // The position of the worker process within the pipeline stages
}

// NewWorkerOptions parses the command-line flags for starting a new worker process and stores them in an
// instance of WorkerOptions
func NewWorkerOptions() *WorkerOptions {
	options := new(WorkerOptions)
	flag.StringVar(&options.MasterAddress, "address", "127.0.0.1",
		"The internet address of the node running the master process")
	flag.IntVar(&options.MasterPort, "config", 8374, "The port number of the master process")
	flag.IntVar(&options.Position, "position", 0, "The position of the worker process within the pipeline stages")
	return options
}
