package worker

import (
	"encoding/gob"
	"net"
	"strconv"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

var waitingForStartPipelineMessage = true

// runFirstStage runs the function of a worker running the first stage
func runFirstStage(listener net.Listener, functionList []types.AnyFunc, myID string, registerType interface{}) {
	go receiveMessages(listener)
	for waitingForStartPipelineMessage {
		// Busy wait lol
	}
	for {
		gob.Register(registerType)
		message := executeStage(functionList, 0, myID, nil)
		encoder := connections.Select()
		if err := encoder.Encode(message); err != nil {
			logMessage(err.Error())
			break
		}
		logPrint("Sent computation results to next stage")
	}
}

// receiveMessage receives messages from the listener
func receiveMessages(listener net.Listener) {
	for {
		message := new(types.Message)
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		decoder := gob.NewDecoder(connection)
		decoder.Decode(message)
		if message.Description == common.MsgAddNextStageAddr {
			nextNodeAddress := (message.Contents).(string)
			connections.AddConnection(nextNodeAddress)
			logPrint("Received next node address")
			continue
		}
		if message.Description == common.MsgStartWorker {
			waitingForStartPipelineMessage = false
			logPrint("Received start pipeline message")
			continue
		}
		logMessage("Received invalid message from " + message.Sender + " of type: " + strconv.Itoa(message.Description))
		connection.Close()
	}
}
