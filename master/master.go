// Package master contains the logic for running the master process of the gopipeline library
package master

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/scheduler"
	"github.com/ffrankies/gopipeline/types"
)

// startListener creates and starts a listener that listens for connections from workers. For each connection, it
// starts a goroutine that reads the messages from the connection.
func startListener(schedule *scheduler.Schedule) (masterAddress string, err error) {
	masterHost := common.GetOutboundIPAddressHack()
	listener, err := net.Listen("tcp", masterHost+":0")
	if err != nil {
		return
	}
	go receiveConnectionsGoRoutine(schedule, listener)
	masterPort := common.GetPortNumberFromListener(listener)
	masterAddress = masterHost + ":" + masterPort
	return
}

// receiveConnectionsGoRoutine is a goroutine that accepts connections from the workers and parses the messages
// received from the workers in separate gosubroutines.
func receiveConnectionsGoRoutine(schedule *scheduler.Schedule, listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnectionFromWorker(schedule, connection)
	}
}

// handleConnectionFromWorker, at the moment, assumes no further communication from the worker node. Thus, it assumes
// the worker sends it's listener address, parses the message as such, and then closes the connection.
// In the future, it should check message type, and either do the above, or update a stage/node's statistics
func handleConnectionFromWorker(schedule *scheduler.Schedule, connection net.Conn) {
	gob.Register(types.MessageStageInfo{})
	decoder := gob.NewDecoder(connection)
	message := new(types.Message)
	decoder.Decode(message)
	if message.Description == common.MsgStageInfo {
		fmt.Println("Type of contents = ", reflect.TypeOf(message.Contents))
		stageInfo, ok := (message.Contents).(types.MessageStageInfo) //.(types.MessageStageInfo)
		if ok {
			schedule.StageList.Find(message.Sender).NetAddress = stageInfo.Address
			schedule.StageList.Find(message.Sender).PID = stageInfo.PID
		} else {
			fmt.Println("Not OK!")
		}
	} else {
		fmt.Println("Received invalid message type from", message.Sender)
	}
	connection.Close()
}

// buildWorkerCommand builds the command with which to start a worker
func buildWorkerCommand(program string, masterAddress string, stageID string, position int) string {
	command := program + " -address=" + masterAddress
	command += " -id=" + stageID
	command += " -position=" + strconv.Itoa(position)
	command += " worker"
	return command
}

// establishInitialWorkerCommunication establishes initial communication between workers by telling them the address
// of the next worker in the pipeline
func establishWorkerCommunication(schedule *scheduler.Schedule, numPositions int) {
	for position := 1; position < numPositions; position++ {
		nextWorker := schedule.StageList.FindByPosition(position)
		currentWorker := schedule.StageList.FindByPosition(position - 1)
		sendNextWorkerAddress(currentWorker, nextWorker)
	}
}

// sendNextWorkerAddress sends the next worker's address to the given worker
func sendNextWorkerAddress(currentWorker *types.PipelineStage, nextWorker *types.PipelineStage) {
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

// startWorkers starts the worker at position 0, thereby kick-starting the pipeline
func startWorkers(schedule *scheduler.Schedule) {
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgStartWorker
	firstStage := schedule.StageList.FindByPosition(0)
	connection, err := net.Dial("tcp", firstStage.NetAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connection)
	encoder.Encode(message)
	fmt.Println("Started stage:", firstStage.StageID)
}

// setUpSignalHandler sets up a signal handler for clean exit on termination
func setUpSignalHandler(schedule *scheduler.Schedule, config *Config) {
	signalHandlerChannel := make(chan os.Signal, 1)
	signal.Notify(signalHandlerChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			receivedSignal := <-signalHandlerChannel
			fmt.Println("Received signal:", receivedSignal)
			fmt.Println("Performing cleanup...")
			schedule.StageList.WaitUntilAllListenerPortsUpdated()
			for _, stage := range schedule.StageList.List {
				sshConnection := NewSSHConnection(stage.Host, config.SSHUser, config.SSHPort)
				command := "kill " + strconv.Itoa(stage.PID)
				sshConnection.RunCommand(command)
			}
			os.Exit(0)
		}
	}()
}

// Run executes the main logic of the "master" node.
// This involves setting up the pipeline stages, and starting worker processes on each node in the pipeline.
func Run(options *common.MasterOptions, functionList []types.AnyFunc) {
	schedule := scheduler.NewSchedule()
	config := NewConfig(options.ConfigPath)
	setUpSignalHandler(schedule, config)
	fmt.Println("=====Doing initial scheduling=====")
	schedule.Static(functionList, config.NodeList)
	masterAddress, err := startListener(schedule)
	if err != nil {
		panic(err)
	}
	fmt.Println("=====Starting workers=====")
	for _, stage := range schedule.StageList.List {
		sshConnection := NewSSHConnection(stage.Host, config.SSHUser, config.SSHPort)
		command := buildWorkerCommand(options.Program, masterAddress, stage.StageID, stage.Position)
		fmt.Println("Running command:", command, "on node:", stage.Host)
		go sshConnection.RunCommand(command)
	}
	fmt.Println("=====Waiting for workers to send their net addresses=====")
	schedule.StageList.WaitUntilAllListenerPortsUpdated()
	fmt.Println("=====Setting up communication between workers=====")
	establishWorkerCommunication(schedule, len(functionList))
	startWorkers(schedule)
	scheduler.Dynamic()
}
