package worker

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// StageID is the ID of this worker
var StageID string

func logMessage(message string) {
	message = "Worker " + StageID + ": " + message
	fmt.Println(message)
}

// sendInfoToMaster opens a connection to the master node, and sends the address of its listener and the pid of this
// stage's worker process
func sendInfoToMaster(masterAddress string, myID string, myAddress string) {
	message := new(types.Message)
	message.Sender = myID
	message.Description = common.MsgStageInfo
	stageInfo := types.MessageStageInfo{Address: myAddress, PID: os.Getpid()}
	message.Contents = stageInfo
	connection, err := net.Dial("tcp", masterAddress)
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	gob.Register(types.MessageStageInfo{})
	encoder := gob.NewEncoder(connection)
	err = encoder.Encode(message)
	if err != nil {
		panic(err)
	}
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
	}
	logMessage("Received invalid message from " + message.Sender + " of type: " + strconv.Itoa(message.Description))
	return ""
}

// runFirstStage runs the function of a worker running the first stage
func runFirstStage(nextNodeAddress string, functionList []types.AnyFunc, myID string) {
	fmt.Println("Running first stage")
	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		fmt.Println("Starting computation...")
		message := new(types.Message)
		result := functionList[0]()
		message.Sender = myID
		message.Description = common.MsgStageResult
		message.Contents = result
		fmt.Println("Sending results...")
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
		logMessage("Starting Pipeline")
	} else {
		logMessage("Received invalid message from: " + message.Sender + " Expected: MsgStartWorker, and instead " +
			" received " + strconv.Itoa(message.Description))
	}
}

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc, registerType interface{}) {
	gob.Register(registerType)
	StageID = options.StageID

	// Listens for both the master and any other connection
	myAddress := common.GetOutboundIPAddressHack()
	listener, err := net.Listen("tcp", myAddress+":0")
	if err != err {
		panic(err)
	}

	// Sends my address as a struct data to the master.
	myPortNumber := common.GetPortNumberFromListener(listener)
	myNetAddress := common.CombineAddressAndPort(myAddress, myPortNumber)
	sendInfoToMaster(options.MasterAddress, options.StageID, myNetAddress)
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
