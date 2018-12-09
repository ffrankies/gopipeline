package worker

import (
	"encoding/gob"
	"strconv"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// decodeInput decodes input from a previous stage
func decodeInput(decoder *gob.Decoder, registerType interface{}) (input interface{}, messageDesc int, err error) {
	//que := makeQueue(10) //check the size of the queue

	gob.Register(registerType)
	message := new(types.Message)
	err = decoder.Decode(message)
	if err != nil {
		logMessage(err.Error())
	}
	input = message.Contents
	messageDesc = message.Description
	//que.Push(&Element{input}) //FILL THE PUSH PART OF THE QUEUE
	return
}

// executeAndSend computes the result of the stage and sends it to the next stage.
func executeAndSend(functionList []types.AnyFunc, position int, myID string, inputQueue *Queue, outputQueue *Queue) {
	go send(outputQueue)
	for {
		input := inputQueue.Pop()
		message := executeStage(functionList, position, myID, input)
		outputQueue.Push(message)
		logPrint("Finished execution")
	}
}

// send sends results from the output queue to the next node
func send(outputQueue *Queue) {
	for {
		output := outputQueue.Pop()
		encoder := connections.Select()
		if err := encoder.Encode(output); err != nil {
			logMessage(err.Error())
			break
		}
		logPrint("Sent computation results to next stage")
	}
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

// executeOnly computes the result of the stage and logs the time at which the computation completed.
func executeOnly(functionList []types.AnyFunc, position int, myID string, inputQueue *Queue, outputQueue *Queue) {
	for {
		input := inputQueue.Pop()
		logPerformance(common.PerfStartExec)
		executeStage(functionList, position, myID, input)
		message := finishedExecutionMessage()
		outputQueue.Push(message)
		logPerformance(common.PerfEndExec)
	}
}

// finishedExecutionMessage creates a message indicating that execution has finished
func finishedExecutionMessage() *types.Message {
	currentTime := time.Now().UnixNano()
	message := new(types.Message)
	message.Sender = StageID
	message.Description = common.MsgEndExecution
	message.Contents = strconv.FormatInt(currentTime, 10)
	return message
}
