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

// The list of pipeline stages
var pipelineStageList = NewPipelineStageList()

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
	numFunctions := len(functionList) - pipelineStageList.Length()
	numNodes := len(nodeList) - len(pipelineNodeList)
	density := math.Ceil(float64(numFunctions) / float64(numNodes))
	return int(density)
}

// assignStageToNode assigns a single pipeline stage (function) to a single node
func assignStageToNode(function types.AnyFunc, nodeAddress string) {
	pipelineNode, foundInList := findNode(nodeAddress)
	pipelineStage := pipelineStageList.AddStage(nodeAddress, pipelineStageList.Length())
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		pipelineNodeList = append(pipelineNodeList, pipelineNode)
	}
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

// handleConnectionFromWorker, at the moment, assumes no further communication from the worker node. Thus, it assumes
// the worker sends it's listener address, parses the message as such, and then closes the connection.
// In the future, it should check message type, and either do the above, or update a stage/node's statistics
func handleConnectionFromWorker(connection net.Conn) {
	decoder := gob.NewDecoder(connection)
	message := new(types.Message)
	decoder.Decode(message)
	if message.Description == common.MsgStageAddr {
		nextNodeAddress := (message.Contents).(string)
		pipelineStageList.Find(message.Sender).NetAddress = nextNodeAddress
	} else {
		fmt.Println("Received invalid message type from", message.Sender)
	}
	connection.Close()
}

// buildWorkerCommand builds the command with which to start a worker
func buildWorkerCommand(program string, masterAddress string, stageID string) string {
	command := program + " -address=" + masterAddress + " -id=" + stageID + " worker"
	return command
}

// establishInitialWorkerCommunication establishes initial communication between workers by telling them the address
// of the next worker in the pipeline
func establishWorkerCommunication(numPositions int) {
	for position := 1; position < numPositions; position++ {
		nextWorker := pipelineStageList.FindByPosition(position)
		currentWorker := pipelineStageList.FindByPosition(position - 1)
		sendNextWorkerAddress(currentWorker, nextWorker)
	}
}

// sendNextWorkerAddress sends the next worker's address to the given worker
func sendNextWorkerAddress(currentWorker *PipelineStage, nextWorker *PipelineStage) {
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgNextStageAddr
	message.Contents = nextWorker.NetAddress
	fmt.Println("Setting up connection to:", currentWorker.NetAddress)
	connection, err := net.Dial("tcp", currentWorker.NetAddress)
	// defer connection.Close()
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
	fmt.Println("Sent addr", nextWorker.NetAddress, "to:", currentWorker.NetAddress)
}

func startWorkers() {
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgStartWorker
	firstStage := pipelineStageList.FindByPosition(0)
	connection, err := net.Dial("tcp", firstStage.NetAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
	fmt.Println("Started stage:", firstStage.StageID)
}

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
func Run(options *common.MasterOptions, functionList []types.AnyFunc) {
	config := NewConfig(options.ConfigPath)
	fmt.Println("=====Doing initial scheduling=====")
	matchStagesToNodes(functionList, config.NodeList)
	masterAddress, err := startListener()
	if err != nil {
		panic(err)
	}
	fmt.Println("=====Starting workers=====")
	for _, stage := range pipelineStageList.List {
		sshConnection := NewSSHConnection(stage.Host, config.SSHUser, config.SSHPort)
		command := buildWorkerCommand(options.Program, masterAddress, stage.StageID)
		fmt.Println("Running command:", command, "on node:", stage.Host)
		go sshConnection.RunCommand(command)
	}
	fmt.Println("=====Waiting for workers to send their net addresses=====")
	pipelineStageList.WaitUntilAllListenerPortsUpdated()
	fmt.Println("=====Setting up communication between workers=====")
	establishWorkerCommunication(len(functionList))
	startWorkers()
	// TODO(): Profit
	for { // Scheduler should work in here
		// Busy wait
	}
}
