package worker

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"

	"github.com/ffrankies/gopipeline/types"
)

// runLastStage runs the function of a worker running the last stage
func runLastStage(listener net.Listener, functionList []types.AnyFunc, myID string, registerType interface{},
	masterAddress string) {

	inputQueue := makeQueue()
	outputQueue := makeQueue()
	go executeOnly(functionList, len(functionList)-1, myID, inputQueue, outputQueue)
	go sendCompletionMessagesToMaster(outputQueue, masterAddress)
	setUpSignalHandler(inputQueue, outputQueue, masterAddress)
	for {
		connectionFromPreviousWorker, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		decoder := gob.NewDecoder(connectionFromPreviousWorker)
		for {
			input, messageDesc, err := decodeInput(decoder, registerType)
			if err != nil {
				break
			}
			if messageDesc == common.MsgStageResult {
				inputQueue.Push(input)
				WorkerStatistics.UpdateBacklog(inputQueue.GetLength())
			} else {
				logMessage("ERROR: Last stage received unexpected message: " + string(messageDesc))
			}
		}
	}
}

// sendCompletionMessagesToMaster sends messages indicating complication of a run through the pipeline
func sendCompletionMessagesToMaster(outputQueue *Queue, masterAddress string) {
	for {
		message := outputQueue.Pop()
		connectionToMaster, err := net.Dial("tcp", masterAddress)
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
		connectionToMaster.Close()
	}
}
