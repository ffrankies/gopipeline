package worker

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// sendAddressToMaster opens a connection to the master node, and sends the address of its listener
func sendAddressToMaster(masterAddress string, myID string, myAddress string) {
	message := new(types.Message)
	message.Sender = myID
	message.Description = common.MsgStageAddr
	message.Contents = myAddress
	connection, err := net.Dial("tcp", masterAddress)
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
}

// receiveAddressOfNextNode listens for a message on the listener, assumes it is from master and contains the address
// of the next code, and parses it as such
func receiveAddressOfNextNode(listener net.Listener) string {
	message := new(types.Message)
	connection, err := listener.Accept()
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(connection)
	decoder.Decode(message)
	if message.Description == common.MsgNextStageAddr {
		nextNodeAddress := (message.Contents).(string)
		return nextNodeAddress
	} else {
		fmt.Println("Worker: Received invalid message from:", message.Sender)
		return ""
	}
}

// runFirstStage runs the function of a worker running the first stage
func runFirstStage(nextNodeAddress string, functionList []types.AnyFunc, myID string) {
	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		message := new(types.Message)
		result := functionList[0](message.Contents)
		message.Sender = myID
		message.Description = common.MsgStageResult
		message.Contents = result
		encoder.Encode(message)
	}
}

// runLastStage runs the function of a worker running the last stage
func runLastStage(listener net.Listener, functionList []types.AnyFunc) {
	connectionFromPreviousWorker, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(connectionFromPreviousWorker)
	for {
		message := new(types.Message)
		decoder.Decode(message)
		functionList[len(functionList)-1](message.Contents)
	}
}

// runIntermediateStage runs the function of a worker running an intermediate stage
func runIntermediateStage(listener net.Listener, nextNodeAddress string, functionList []types.AnyFunc, myID string,
	position int) {
	connectionFromPreviousWorker, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	decoder := gob.NewDecoder(connectionFromPreviousWorker)
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		message := new(types.Message)
		decoder.Decode(message)
		result := functionList[position](message.Contents)
		message.Sender = myID
		message.Description = common.MsgStageResult
		message.Contents = result
		encoder.Encode(message)
	}
}

// waitForStartCommand tells a worker to wait for
func waitForStartCommand(listener net.Listener) {
	message := new(types.Message)
	connection, err := listener.Accept()
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(connection)
	decoder.Decode(message)
	if message.Description == common.MsgStartWorker {
		fmt.Println("Starting pipeline")
	} else {
		fmt.Println("Worker: Received invalid message from:", message.Sender, "Expected: MsgStartWorker, and instead",
			"received", message.Description)
	}
}

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc) {

	// Listens for both the master and any other connection
	myAddress := common.GetOutboundIPAddressHack()
	listener, err := net.Listen("tcp", myAddress+":0")
	if err != err {
		panic(err)
	}

	// Sends my address as a struct data to the master.
	myPortNumber := common.GetPortNumberFromListener(listener)
	myNetAddress := common.CombineAddressAndPort(myAddress, myPortNumber)
	sendAddressToMaster(options.MasterAddress, options.StageID, myNetAddress)
	isLastStage := options.Position == len(functionList)-1
	var nextNodeAddress string
	if !isLastStage {
		nextNodeAddress = receiveAddressOfNextNode(listener)
	}
	// Get data from previous worker, process it, and send results to the next worker
	if options.Position == 0 {
		waitForStartCommand(listener)
		runFirstStage(nextNodeAddress, functionList, options.StageID)
	}
	if isLastStage {
		runLastStage(listener, functionList)
	}
	runIntermediateStage(listener, nextNodeAddress, functionList, options.StageID, options.Position)
}
