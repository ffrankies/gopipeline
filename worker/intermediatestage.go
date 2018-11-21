package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/types"
)

// Note: erase later. msg receiving and result reciving should be on different portions,
//and it should lock one before executing the other.
//So that they donot run simultaiously.

// runIntermediateStage runs the function of a worker running an intermediate stage
func runIntermediateStage(listener net.Listener, nextNodeAddress string, functionList []types.AnyFunc, myID string,
	position int, registerType interface{}) {
	logPrint("In run Intermediate Stage module")
	queue := makeQueue()
	go exeecuteAndSend(functionList, position, myID, queue, nextNodeAddress)
	for {
		connectionFromPreviousWorker, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		decoder := gob.NewDecoder(connectionFromPreviousWorker)

		for {
			logMessage("Starting intermediate computation...")
			input, err := decodeInput(decoder, registerType)
			if err != nil {
				break
			}
			queue.Push(input)

		}
	}
}
