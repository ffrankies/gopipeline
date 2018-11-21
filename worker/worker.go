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

// WorkerStatistics is the performance statistics of this worker process
var WorkerStatistics = new(types.WorkerStats)

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
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	gob.Register(types.MessageStageInfo{})
	encoder := gob.NewEncoder(connection)
	err = encoder.Encode(message)
	if err != nil {
		panic(err)
	}
}

// runStage chooses the correct stage function to run, and runs it
func runStage(options *common.WorkerOptions, functionList []types.AnyFunc, listener net.Listener,
	registerType interface{}) {

	isLastStage := options.Position == len(functionList)-1
	var nextNodeAddress string
	if !isLastStage {
		nextNodeAddress = receiveAddressOfNextNode(listener)
	}
	// Get data from previous worker, process it, and send results to the next worker
	if options.Position == 0 {
		waitForStartCommand(listener)
		runFirstStage(nextNodeAddress, functionList, options.StageID, registerType)
	}
	if isLastStage {
		runLastStage(listener, functionList, registerType)
	}
	runIntermediateStage(listener, nextNodeAddress, functionList, options.StageID, options.Position, registerType)
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

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc, registerType interface{}) {
	StageID = options.StageID

	go trackStatsGoroutine(options.MasterAddress, options.StageID)

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
	runStage(options, functionList, listener, registerType)
}
