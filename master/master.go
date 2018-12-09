// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/scheduler"
	"github.com/ffrankies/gopipeline/types"
)

// NumExecutions is the number of times the pipeline has executed from start to finish
var NumExecutions = 0

// NumExecutionsMutex is the mutex for changing NumExecutions
var NumExecutionsMutex = &sync.Mutex{}

// startListener creates and starts a listener that listens for connections from workers. For each connection, it
// starts a goroutine that reads the messages from the connection.
func startListener(schedule *scheduler.Schedule, config *Config) (masterAddress string, err error) {
	masterHost := common.GetOutboundIPAddressHack()
	listener, err := net.Listen("tcp", masterHost+":0")
	if err != nil {
		return
	}
	go receiveConnectionsGoRoutine(schedule, listener, config)
	masterPort := common.GetPortNumberFromListener(listener)
	masterAddress = masterHost + ":" + masterPort
	return
}

// receiveConnectionsGoRoutine is a goroutine that accepts connections from the workers and parses the messages
// received from the workers in separate gosubroutines.
func receiveConnectionsGoRoutine(schedule *scheduler.Schedule, listener net.Listener, config *Config) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnectionFromWorker(schedule, connection, config)
	}
}

// handleConnectionFromWorker, at the moment, assumes no further communication from the worker node. Thus, it assumes
// the worker sends it's listener address, parses the message as such, and then closes the connection.
// In the future, it should check message type, and either do the above, or update a stage/node's statistics
func handleConnectionFromWorker(schedule *scheduler.Schedule, connection net.Conn, config *Config) {
	gob.Register(&types.WorkerStats{})
	gob.Register(types.MessageStageInfo{})
	decoder := gob.NewDecoder(connection)
	message := new(types.Message)
	decoder.Decode(message)
	if message.Description == common.MsgStageInfo {
		schedule.UpdateStageInfo(message)
	} else if message.Description == common.MsgStageStats {
		schedule.UpdateStageStats(message)
	} else if message.Description == common.MsgNotifyExit {
		exitingWorkerID := message.Sender
		schedule.StageList.RemoveWorker(exitingWorkerID)
		schedule.NodeList.RemoveWorker(exitingWorkerID)
	} else if message.Description == common.MsgStartFirstStage {
		logPrint(common.PerfStartWorker + " " + (message.Contents).(string))
	} else if message.Description == common.MsgEndExecution {
		NumExecutionsMutex.Lock()
		NumExecutions++
		if NumExecutions == 1 || NumExecutions == 5 {
			logPrint(common.PerfEndExec + " " + (message.Contents).(string))
		}
		if NumExecutions == 5 {
			NumExecutionsMutex.Unlock()
			killAllWorkersAndExit(schedule, config)
		} else {
			NumExecutionsMutex.Unlock()
		}
	} else {
		fmt.Println("Received invalid message type from", message.Sender)
	}
	connection.Close()
}

// startWorkers starts the worker at position 0, thereby kick-starting the pipeline
func startWorkers(schedule *scheduler.Schedule) {
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgStartWorker
	firstWorker := schedule.StageList.FindWorker("1")
	connection, err := net.Dial("tcp", firstWorker.Address)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
	fmt.Println("Started worker:", firstWorker.ID)
}

// setUpSignalHandler sets up a signal handler for clean exit on termination
func setUpSignalHandler(schedule *scheduler.Schedule, config *Config) {
	signalHandlerChannel := make(chan os.Signal, 1)
	signal.Notify(signalHandlerChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			receivedSignal := <-signalHandlerChannel
			fmt.Println("Received signal:", receivedSignal)
			killAllWorkersAndExit(schedule, config)
		}
	}()
}

// killAllWorkersAndExit kills all workers, then exits master
func killAllWorkersAndExit(schedule *scheduler.Schedule, config *Config) {
	fmt.Println("Performing cleanup...")
	schedule.StageList.WaitUntilAllListenerPortsUpdated()
	for _, stage := range schedule.StageList.List {
		for _, worker := range stage.Workers {
			sshConnection := types.NewSSHConnection(worker.Host, config.SSHUser, config.SSHPort)
			command := "kill " + strconv.Itoa(worker.PID)
			sshConnection.RunCommand(command, nil, nil)
		}
	}
	os.Exit(0)
}

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
func Run(options *common.MasterOptions, functionList []types.AnyFunc) {
	config := NewConfig(options.ConfigPath)
	schedule := scheduler.NewSchedule(
		config.NodeList, config.SSHUser, config.SSHPort, config.UserPath, len(functionList))
	setUpSignalHandler(schedule, config)
	schedule.Static(functionList)
	masterAddress, err := startListener(schedule, config)
	if err != nil {
		panic(err)
	}
	schedule.StartStages(options.Program, masterAddress)
	fmt.Println("=====Waiting for workers to send their net addresses=====")
	schedule.StageList.WaitUntilAllListenerPortsUpdated()
	fmt.Println("=====Setting up communication between workers=====")
	schedule.EstablishWorkerCommunication()
	startWorkers(schedule)
	schedule.Dynamic(options.Program, masterAddress)
}
