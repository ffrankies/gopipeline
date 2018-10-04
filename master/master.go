// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"fmt"

	"github.com/ffrankies/gopipeline"
)

// The list of pipeline stages
var pipelineStageList []PipelineStage

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
// The command is the command to be used to start the worker process.
// The configPath is the path to the config file that contains the login information and node list.
// The functionList is the list of functions to pipeline.
func Run(command string, configPath string, functionList []gopipeline.AnyFunc) {
	stage := NewPipelineStage("127.0.0.1", 10, 11, 12, 0)
	fmt.Println("Stage =", stage)
	config := NewConfig(configPath)
	fmt.Println("Config = ", config)
	for _, address := range config.NodeList {
		sshConnection := NewSSHConnection(address, config.SSHUser, config.SSHPort)
		out, err := sshConnection.RunCommand("hostname")
		if err != nil {
			panic("Failed to run ls: " + err.Error())
		}
		fmt.Println("The hostname is...", out)
		sshConnection.Close()
	}
	// out, err = sshConnection.RunCommand("ps")
	// if err != nil {
	// 	panic("Failed to run ps: " + err.Error())
	// }
	// fmt.Println(out)
}
