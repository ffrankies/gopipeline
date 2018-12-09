package worker

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"sync"
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

// logging mutex
var logMutex = &sync.Mutex{}

// userHomeDir returns the current user's home directory
func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

// opens a log file in the user's home directory
func openLogFile() (fp *os.File) {
	userPath := userHomeDir()
	filePath := userPath + "/gopipeline" + StageNumber + "." + StageID + ".log"
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// logPrint prints a message to the log file
func logPrint(message string) {
	logMutex.Lock()
	f := openLogFile()
	defer f.Close()
	log.SetOutput(f)
	message = "Worker " + StageID + " | Stage " + StageNumber + ": " + message
	log.Println(message)
	logMutex.Unlock()
}

// logMessage prints a message to the console AND the log file
func logMessage(message string) {
	logPrint(message)
	message = "Worker " + StageID + ": " + message
	fmt.Println(message)
}

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

// Run the worker routine
func Run(options *common.WorkerOptions, functionList []types.AnyFunc, registerType interface{}) {
	StageID = options.StageID
	StageNumber = strconv.Itoa(options.Position)
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
