package worker

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// StageID is the ID of this worker
var StageID string

// StageNumber the number of this stage
var StageNumber string

// WorkerStatistics is the performance statistics of this worker process
var WorkerStatistics = new(types.WorkerStats)

// connections is the list of connections to the next nodes
var connections = NewConnections()

// sendInfoToMaster opens a connection to the master node, and sends the address of its listener and the pid of this
// stage's worker process
func sendInfoToMaster(masterAddress string, myID string, myAddress string) {

	logPrint("Sending info to master")
	message := new(types.Message)
	message.Sender = myID
	message.Description = common.MsgStageInfo
	stageInfo := types.MessageStageInfo{Address: myAddress, PID: os.Getpid()}
	message.Contents = stageInfo
	connection, err := net.DialTimeout("tcp", masterAddress, 2*time.Second)
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
	// Get data from previous worker, process it, and send results to the next worker
	logPrint("My position is " + strconv.Itoa(options.Position))
	if options.Position == 0 {
		// waitForStartCommand(listener)
		runFirstStage(listener, functionList, options.StageID, registerType, options.MasterAddress)
	} else if isLastStage {
		runLastStage(listener, functionList, options.StageID, registerType, options.MasterAddress)
	} else {
		runIntermediateStage(listener, functionList, options.StageID, options.Position, registerType, options.MasterAddress)
	}
}

// sendStartMessageToMaster sends messages indicating complication of a run through the pipeline
func sendStartMessageToMaster(masterAddress string) {
	currentTime := time.Now().UnixNano()
	message := new(types.Message)
	message.Sender = StageID
	message.Description = common.MsgStartFirstStage
	message.Contents = strconv.FormatInt(currentTime, 10)
	connectionToMaster, err := net.Dial("tcp", masterAddress)
	defer connectionToMaster.Close()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToMaster)
	err = encoder.Encode(message)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc, registerType interface{}) {
	StageID = options.StageID
	StageNumber = strconv.Itoa(options.Position)
	if StageID == "1" || StageNumber == strconv.Itoa(len(functionList)-1) {
		setupPerformanceLogFile()
		if StageID == "1" {
			sendStartMessageToMaster(options.MasterAddress)
			logPerformance(common.PerfStartWorker)
		}
	}
	setupLogFile()
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
