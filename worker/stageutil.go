package worker

import (
	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// executeStage executes the function this stage is responsible for, and returns the result as a message
func executeStage(functionList []types.AnyFunc, position int, stageID string, input interface{}) *types.Message {
	message := new(types.Message)
	var result interface{}
	if input == nil {
		result = functionList[position]()
	} else {
		result = functionList[position](input)
	}
	message.Sender = stageID
	message.Description = common.MsgStageResult
	message.Contents = result
	return message
}
