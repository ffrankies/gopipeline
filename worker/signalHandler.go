package worker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// setUpSignalHandler sets up a signal handler for clean exit on termination
func setUpSignalHandler(inputQueue *Queue, outputQueue *Queue) {
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
				os.Exit(0)
			}
		}
	}()
}
