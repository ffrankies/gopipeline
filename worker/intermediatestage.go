package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"

	"github.com/ffrankies/gopipeline/types"
)

// runIntermediateStage runs the function of a worker running an intermediate stage
func runIntermediateStage(listener net.Listener, functionList []types.AnyFunc, myID string, position int,
	registerType interface{}, masterAddress string) {

	inputQueue := makeQueue()
	outputQueue := makeQueue()
	go executeAndSend(functionList, position, myID, inputQueue, outputQueue)
	setUpSignalHandler(inputQueue, outputQueue, masterAddress)
	for {
		logPrint("Waiting for connection from whoever")
		listenerConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(listenerConnection, registerType, inputQueue)
	}
}

// handleConnection handles a connection from either previous worker or master
func handleConnection(connection net.Conn, registerType interface{}, inputQueue *Queue) {
	decoder := gob.NewDecoder(connection)
	for {
		input, messageDesc, err := decodeInput(decoder, registerType)
		if err != nil {
			break
		}
		if messageDesc == common.MsgStageResult {
			inputQueue.Push(input)
			WorkerStatistics.UpdateBacklog(inputQueue.GetLength())
			logPrint("Received input from previous worker")
		}
		if messageDesc == common.MsgAddNextStageAddr {
			connections.AddConnection(input.(string))
			logPrint("Received new address from master")
		}
		if messageDesc == common.MsgBreakConnection {
			addressToRemove := input.(string)
			connections.RemoveConnection(addressToRemove)
			logPrint("Removed the worker from the list of connections")
			continue
		}
	}
}
