package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/types"
)

// runIntermediateStage runs the function of a worker running an intermediate stage
func runIntermediateStage(listener net.Listener, nextNodeAddress string, functionList []types.AnyFunc, myID string,
	position int, registerType interface{}) {
	for {
		connectionFromPreviousWorker, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
		decoder := gob.NewDecoder(connectionFromPreviousWorker)
		encoder := gob.NewEncoder(connectionToNextWorker)
		for {
			logMessage("Starting intermediate computation...")
			gob.Register(registerType)
			message := new(types.Message)
			if err := decoder.Decode(message); err != nil {
				logMessage(err.Error())
				break
			}
			message = executeStage(functionList, position, myID, message.Contents)
			if err := encoder.Encode(message); err != nil {
				logMessage(err.Error())
				break
			}
			logMessage("Ending intermediate computation...")
		}
	}
}
