package worker

import (
	"encoding/gob"
	"fmt"
	"log"
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

// log writes log details to a logfile
func logPrint(message string) {
	f, err := os.OpenFile("/Users/bipashabanerjee/go/src/github.com/ffrankies/gopipeline/logFile", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	//test case
	log.Println(message)
}

// sendInfoToMaster opens a connection to the master node, and sends the address of its listener and the pid of this
// stage's worker process
func sendInfoToMaster(masterAddress string, myID string, myAddress string) {
	logPrint("In  send Info to Message block")
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
	logPrint("In receive anddress of Next Node block")
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
	logPrint("In run first stage block")
	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		message := new(types.Message)
		result := functionList[0]()
		message.Sender = myID
		message.Description = common.MsgStageResult
		message.Contents = result
		encoder.Encode(message)
	}
}

// runLastStage runs the function of a worker running the last stage
func runLastStage(listener net.Listener, functionList []types.AnyFunc) {
	logPrint("run the last stage block")
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
	logPrint("run intermediate stage block")
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

	logPrint("Waiting for the start command")
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
		log.Println("Starting Pipeline")
	} else {
		logMessage("Received invalid message from: " + message.Sender + " Expected: MsgStartWorker, and instead " +
			" received " + strconv.Itoa(message.Description))
	}
}

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc) {
	//create your file with desired read/write permissions
	logPrint("Run Worker Called")
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
