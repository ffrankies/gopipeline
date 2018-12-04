package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/types"
)

// runIntermediateStage runs the function of a worker running an intermediate stage
func runIntermediateStage(listener net.Listener, nextNodeAddress string, functionList []types.AnyFunc, myID string,
	position int, registerType interface{}) {
	logPrint("In run Intermediate Stage module")
	queue := makeQueue()
	go executeAndSend(functionList, position, myID, queue, nextNodeAddress)
	for {
		connectionFromPreviousWorker, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		decoder := gob.NewDecoder(connectionFromPreviousWorker)
		for {
			input, err := decodeInput(decoder, registerType)
			if err != nil {
				break
			}
			queue.Push(input)
		}
	}
}
