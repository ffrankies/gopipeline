package worker

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// setUpSignalHandler sets up a signal handler for clean exit on termination
func setUpSignalHandler(inputQueue *Queue, outputQueue *Queue, masterAddress string) {
	signalHandlerChannel := make(chan os.Signal, 1)
	signal.Notify(signalHandlerChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			receivedSignal := <-signalHandlerChannel
			fmt.Println("Received signal:", receivedSignal)
			if receivedSignal == syscall.SIGINT || receivedSignal == syscall.SIGTERM {
				fmt.Println("Performing cleanup...")
				connections.CloseAll()
				// Figure out how to kill listener
				os.Exit(-1)
			}
			if receivedSignal == syscall.SIGUSR1 {
				fmt.Println("Received a SigUSR1 signal...")
				if inputQueue != nil {
					inputQueue.WaitUntilEmpty()
				}
				if outputQueue != nil {
					outputQueue.WaitUntilEmpty()
				}
				connections.CloseAll()
				// Figure out how to kill listener
				notifyMasterOfExit(masterAddress)
				os.Exit(0)
			}
		}
	}()
}

// notifyMasterOfExit notifies the master that this node is about to exit
func notifyMasterOfExit(masterAddress string) {
	message := new(types.Message)
	message.Sender = StageID
	message.Description = common.MsgNotifyExit
	message.Contents = StageID
	connectionToMaster, err := net.Dial("tcp", masterAddress)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer connectionToMaster.Close()
	encoder := gob.NewEncoder(connectionToMaster)
	err = encoder.Encode(message)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
