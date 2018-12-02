package worker

import (
	"encoding/gob"
	"net"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// decodeInput decodes input from a previous stage
func decodeInput(decoder *gob.Decoder, registerType interface{}) (input interface{}, err error) {
	//que := makeQueue(10) //check the size of the queue

	gob.Register(registerType)
	message := new(types.Message)
	err = decoder.Decode(message)
	if err != nil {
		logMessage(err.Error())
	}
	input = message.Contents
	//que.Push(&Element{input}) //FILL THE PUSH PART OF THE QUEUE
	return
}

// executeStage executes the function this stage is responsible for, and returns the result as a message
func executeStage(functionList []types.AnyFunc, position int, stageID string, input interface{}) *types.Message {
	message := new(types.Message)
	var result interface{}
	timerStart := time.Now()
	if input == nil {
		result = functionList[position](nil)
	} else {
		result = functionList[position](input)
	}
	WorkerStatistics.UpdateExecutionTime(time.Since(timerStart))
	message.Sender = stageID
	message.Description = common.MsgStageResult
	message.Contents = result
	return message
}

// executeAndSend computes the result of the stage and sends it to the next stage.
func executeAndSend(functionList []types.AnyFunc, position int, myID string, queue *Queue, nextNodeAddress string) {

	connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(connectionToNextWorker)
	for {
		logMessage("Starting Computation")
		input := queue.Pop()
		message := executeStage(functionList, position, myID, input)
		if err := encoder.Encode(message); err != nil {
			logMessage(err.Error())
			break
		}

	}

}
