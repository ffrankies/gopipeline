// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"fmt"
	"math"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// The list of pipeline stages
var pipelineStageList []*PipelineStage

// The list of pipeline nodes
var pipelineNodeList []*PipelineNode

// Matches pipeline stages (functions) to the nodes on which they will run.
// The algorithm tries to spread the functions out among the nodes, but if that isn't possible
// (there are more functions than nodes), it will automatically bunch functions together.
func matchStagesToNodes(functionList []types.AnyFunc, nodeList []string) {
	density := calculateFunctionDensity(functionList, nodeList)
	counter := 0
	nodeIndex := 0
	for _, function := range functionList {
		assignStageToNode(function, nodeList[nodeIndex])
		counter++
		if counter == density {
			nodeIndex++
			counter = 0
			density = calculateFunctionDensity(functionList, nodeList)
		}
	}
}

// calculateFunctionDensity calculates the initial function density in the pipeline
func calculateFunctionDensity(functionList []types.AnyFunc, nodeList []string) int {
	numFunctions := len(functionList) - len(pipelineStageList)
	numNodes := len(nodeList) - len(pipelineNodeList)
	density := math.Ceil(float64(numFunctions) / float64(numNodes))
	return int(density)
}

// assignStageToNode assigns a single pipeline stage (function) to a single node
func assignStageToNode(function types.AnyFunc, nodeAddress string) {
	pipelineNode, foundInList := findNode(nodeAddress)
	pipelineStage := NewPipelineStage(nodeAddress, 0, 0, 0, len(pipelineStageList))
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		pipelineNodeList = append(pipelineNodeList, pipelineNode)
	}
	pipelineStageList = append(pipelineStageList, pipelineStage)
}

// findNode finds a particular PipelineNode in the pipelineNodeList. If the Node is not found,
// findNode creates a new PipelineNode
func findNode(nodeAddress string) (pipelineNode *PipelineNode, foundInList bool) {
	for _, node := range pipelineNodeList {
		if node.Address == nodeAddress {
			pipelineNode = node
			foundInList = true
			return
		}
	}
	pipelineNode = NewPipelineNode(nodeAddress, len(pipelineNodeList))
	foundInList = false
	return
}

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
// The command is the command to be used to start the worker process.
// The configPath is the path to the config file that contains the login information and node list.
// The functionList is the list of functions to pipeline.
func Run(options *common.MasterOptions, functionList []types.AnyFunc) {
	fmt.Println("My options are...", options)
	config := NewConfig(options.ConfigPath)
	fmt.Println("Config = ", config)
	matchStagesToNodes(functionList, config.NodeList)
	fmt.Println("=====Node List=====")
	fmt.Println(pipelineNodeList)
	fmt.Println("=====Stage List=====")
	fmt.Println(pipelineStageList)
	// TODO(): Set up a server for communicating with worker processes
	for _, stage := range pipelineStageList {
		sshConnection := NewSSHConnection(stage.NodeAddress, config.SSHUser, config.SSHPort)
		// TODO(): Create command using options.Program, server address and port number, and stage.Position
		out, err := sshConnection.RunCommand(options.Program)
		if err != nil {
			panic("Failed to run, " + options.Program + ": " + err.Error())
		}
		fmt.Println("Stage number,", stage.Position, "is running on host:", out)
		sshConnection.Close()
	}
}
