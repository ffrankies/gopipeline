package worker

import (
	"encoding/gob"
	"net"

	"github.com/ffrankies/gopipeline/types"
)

// runLastStage runs the function of a worker running the last stage
func runLastStage(listener net.Listener, functionList []types.AnyFunc, registerType interface{}) {
	for {
		connectionFromPreviousWorker, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		decoder := gob.NewDecoder(connectionFromPreviousWorker)
		for {
			logMessage("Starting last stage computation...")
			gob.Register(registerType)
			message := new(types.Message)
			if err := decoder.Decode(message); err != nil {
				logMessage(err.Error())
				break
			}
			functionList[len(functionList)-1](message.Contents)
			logMessage("Ending last stage computation...")
		}
	}
}
