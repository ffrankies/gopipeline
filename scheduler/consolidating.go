package scheduler

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// findWorkerToMove searches for the worker that is on the node that is after the given node. It should be using less memory than
// avaiable in the node at the given position.
func (schedule *Schedule) findWorkerToMove(position int, availableMemory uint64) *types.Worker {
	for _, node := range schedule.NodeList.List {
		if position < node.Position {
			for _, worker := range node.Workers {
				if worker.Stats.WorkerMemoryUsage < availableMemory && worker.Stats.ExecutionTime > 0 {
					return worker
				}
			}
		}
	}
	return nil
}

// breakConnection closes the connection between the worker and all the other workers who sends the results to it
func (schedule *Schedule) breakConnection(oldWorkerAddress string, position int) {
	if position == 0 {
		return
	}
	message := new(types.Message)
	message.Sender = "0"
	message.Description = common.MsgBreakConnection
	message.Contents = oldWorkerAddress
	previousStage := schedule.StageList.FindByPosition(position - 1)
	for _, worker := range previousStage.Workers {
		fmt.Println("Setting up connection to:", worker.Address)
		connection, err := net.Dial("tcp", worker.Address)
		if err != nil {
			panic(err)
		}
		encoder := gob.NewEncoder(connection)
		encoder.Encode(message)
		connection.Close()
		fmt.Println("Closed connection between " + worker.ID + " and " + oldWorkerAddress)
	}
}

// flushAndStopWorker sends a signal to the worker to flush its queue and kill itself
func (schedule *Schedule) flushAndStopWorker(worker *types.Worker) {
	sshConnection := types.NewSSHConnection(worker.Host, schedule.sshUser, schedule.sshPort)
	command := "kill -SIGUSR1 " + strconv.Itoa(worker.PID)
	fmt.Println("Running command:", command, "on node:", worker.Host)
	go sshConnection.RunCommand(command, nil, nil)
}

// moveStages moves the data for processing from the current node to the previous node if it
// has memory available for usage
func (schedule *Schedule) moveStages(program string, masterAddress string) {
	fmt.Println("Consolidating workers")
	for _, node := range schedule.NodeList.List {
		availableMemory := node.AvailableMemory()
		worker := schedule.findWorkerToMove(node.Position, availableMemory)
		if worker == nil {
			continue
		}
		fmt.Println("Moving worker " + worker.ID + " to node " + node.Address)
		newWorker := schedule.AssignWorkerToNode(worker.Stage, node)
		schedule.startWorker(newWorker, program, masterAddress)
		fmt.Println("Waiting for the new worker to send info...")
		if err := schedule.waitForWorkerToSendInfo(newWorker); err != nil {
			panic(err)
		}
		fmt.Println("Done waiting for worker to send info...")
		schedule.setUpNewWorkerCommunication(newWorker)
		schedule.breakConnection(worker.Address, worker.Stage)
		schedule.flushAndStopWorker(worker)
		break
	}
}
