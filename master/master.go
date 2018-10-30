// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

type Address struct {
	Data string
}
type Pipeline struct {
	NodeNumber string
}

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
	pipelineStage := NewPipelineStage(nodeAddress, len(pipelineStageList))
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

// startListener creates and starts a listener that listens for connections from workers. For each connection, it
// starts a goroutine that reads the messages from the connection.
func startListener() (masterAddress string, err error) {
	masterHost := common.GetOutboundIPAddressHack()
	listener, err := net.Listen("tcp", masterHost+":0")
	if err != nil {
		return
	}
	go receiveConnectionsGoRoutine(listener)
	masterPort := common.GetPortNumberFromListener(listener)
	masterAddress = masterHost + ":" + masterPort
	return
}

// receiveConnectionsGoRoutine is a goroutine that accepts connections from the workers and parses the messages
// received from the workers in separate gosubroutines.
func receiveConnectionsGoRoutine(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnectionFromWorker(connection)
	}
}

func handleConnectionFromWorker(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	message := new(types.Message)
	dec.Decode(message)
	nextNodeAddress := (message.Contents).(string)
	fmt.Println("next node address:", nextNodeAddress, "gotten from:", message.Sender)
	// conn1, err := net.Dial("tcp", workerAddress)
	// if err != nil {
	// 	panic(err)
	// }
	// // Find the next address..????????????
	// conn1.Write([]byte("next Address"))
	// conn1.Close()
}

// buildWorkerCommand builds the command with which to start a worker
func buildWorkerCommand(program string, masterAddress string) string {
	command := program + " -address=" + masterAddress + " worker"
	return command
}

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
func Run(options *common.MasterOptions, functionList []types.AnyFunc) {
	config := NewConfig(options.ConfigPath)
	matchStagesToNodes(functionList, config.NodeList)
	fmt.Println("=====Node List=====")
	fmt.Println(pipelineNodeList)
	fmt.Println("=====Stage List=====")
	fmt.Println(pipelineStageList)
	masterAddress, err := startListener()
	if err != nil {
		panic(err)
	}
	for _, stage := range pipelineStageList {
		sshConnection := NewSSHConnection(stage.Host, config.SSHUser, config.SSHPort)
		command := buildWorkerCommand(options.Program, masterAddress)
		fmt.Println("Running command:", command, "on node:", stage.Host)
		go sshConnection.RunCommand(command)
	}
	// TODO(): Get each worker's listener port
	// TODO(): Send worker_{i}'s listener port to worker_{i+1}
	// TODO(): Profit
	for { // Scheduler should work in here
		// Busy wait
	}
}
