package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/internal/common"

	"github.com/ffrankies/gopipeline/types"
)

// runLastStage runs the function of a worker running the last stage
func runLastStage(listener net.Listener, functionList []types.AnyFunc, myID string, registerType interface{}, masterAddress string) {
	queue := makeQueue()
	go executeOnly(functionList, len(functionList)-1, myID, queue)
	setUpSignalHandler(nil, queue, masterAddress)
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
				queue.Push(input)
				WorkerStatistics.UpdateBacklog(queue.GetLength())
			} else {
				logMessage("ERROR: Last stage received unexpected message: " + string(messageDesc))
			}
		}
	}
}
