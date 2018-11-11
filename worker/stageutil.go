package worker

import (
	"encoding/gob"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// decodeInput decodes input from a previous stage
func decodeInput(decoder *gob.Decoder, registerType interface{}) (input interface{}, err error) {
	gob.Register(registerType)
	message := new(types.Message)
	err = decoder.Decode(message)
	if err != nil {
		logMessage(err.Error())
	}
	input = message.Contents
	return
}

// executeStage executes the function this stage is responsible for, and returns the result as a message
func executeStage(functionList []types.AnyFunc, position int, stageID string, input interface{}) *types.Message {
	message := new(types.Message)
	var result interface{}
	timerStart := time.Now()
	if input == nil {
		result = functionList[position]()
	} else {
		result = functionList[position](input)
	}
	WorkerStatistics.UpdateExecutionTime(time.Since(timerStart))
	message.Sender = stageID
	message.Description = common.MsgStageResult
	message.Contents = result
	return message
}
