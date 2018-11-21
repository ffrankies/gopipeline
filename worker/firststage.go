package worker

import (
	"encoding/gob"
	"net"
	"strconv"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// waitForStartCommand tells a worker to wait for
func waitForStartCommand(listener net.Listener) {
	logPrint("Waiting for the start command")
	message := new(types.Message)
	connection, err := listener.Accept()
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(connection)
	decoder.Decode(message)
	if message.Description == common.MsgStartWorker {
		logMessage("Starting Pipeline")
	} else {
		logMessage("Received invalid message from: " + message.Sender + " Expected: MsgStartWorker, and instead " +
			" received " + strconv.Itoa(message.Description))
	}
}

// runFirstStage runs the function of a worker running the first stage
func runFirstStage(nextNodeAddress string, functionList []types.AnyFunc, myID string, registerType interface{}) {
	logPrint("In runFirstStage module")
	for {
		connectionToNextWorker, err := net.Dial("tcp", nextNodeAddress)
		if err != nil {
			panic(err)
		}
		encoder := gob.NewEncoder(connectionToNextWorker)
		for {
			logMessage("Starting computation...")
			gob.Register(registerType)
			message := executeStage(functionList, 0, myID, nil)
			err = encoder.Encode(message)
			if err != nil {
				logMessage(err.Error())
				break
			}
			logMessage("Sent results...")
		}
	}
}
