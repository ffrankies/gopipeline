package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// sendPortNumberToMaster opens a connection to the master node, and sends the port number of its listener
func sendPortNumberToMaster(masterAddress string, myID string, myPortNumber string) {
	message := new(types.Message)
	message.Sender = myID
	message.Contents = myPortNumber
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
	nextNodeAddress := (message.Contents).(string)
	return nextNodeAddress
}

// runFirstStage runs the function of a worker running the first stage
func runFirstStage(nextNodeAddress string, functionList []types.AnyFunc, myAddress string) {
	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		message := new(types.Message)
		result := functionList[0](message.Contents)
		message.Sender = myAddress
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
func runIntermediateStage(listener net.Listener, nextNodeAddress string, functionList []types.AnyFunc, myAddress string,
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
		message.Sender = myAddress
		message.Contents = result
		encoder.Encode(message)
	}
}

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc) {

	// Listens for both the master and any other connection
	listener, err := net.Listen("tcp", "localhost:0")
	if err != err {
		panic(err)
	}

	// Sends my address as a struct data to the master.
	myAddress := common.GetOutboundIPAddressHack()
	myPortNumber := common.GetPortNumberFromListener(listener)
	sendPortNumberToMaster(options.MasterAddress, options.StageID, myPortNumber)
	isLastWorker := options.Position == len(functionList)-1
	var nextNodeAddress string
	if !isLastWorker {
		nextNodeAddress = receiveAddressOfNextNode(listener)
	}

	// Get data from previous worker, process it, and send results to the next worker
	if options.Position == 0 {
		runFirstStage(nextNodeAddress, functionList, myAddress)
	}
	if isLastWorker {
		runLastStage(listener, functionList)
	}
	runIntermediateStage(listener, nextNodeAddress, functionList, myAddress, options.Position)
}
