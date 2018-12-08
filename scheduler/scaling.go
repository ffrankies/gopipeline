package scheduler

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ffrankies/gopipeline/internal/common"
	"github.com/ffrankies/gopipeline/types"
)

// scaleStage scales a Bottleneck stage out to a free node
func (schedule *Schedule) scaleStage(position int, numToScale int, program string, masterAddress string) {
	numScaled := 0
	fmt.Println(numScaled, "|", numToScale)
	for numScaled < numToScale {
		if position == -1 {
			return
		}
		// For now, only scale on free nodes
		if schedule.freeNodeList.Length() < 1 {
			return
		}
		newWorker := schedule.AssignWorkerToFreeNode(position)
		schedule.startWorker(newWorker, program, masterAddress)
		fmt.Println("Waiting for worker to send info...")
		if err := schedule.waitForWorkerToSendInfo(newWorker); err != nil {
			panic(err)
		}
		fmt.Println("Done waiting for worker to send info...")
		schedule.setUpNewWorkerCommunication(newWorker)
		numScaled++
	}
	// TODO(): also scale on underutilized nodes
}

// waitForWorkerToSendInfo busy waits until the worker sends its info
func (schedule *Schedule) waitForWorkerToSendInfo(worker *types.Worker) error {
	for worker.PID == -1 {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("PID has been updated")
	if worker.PID == -2 {
		return errors.New("ERROR: Worker could not be started")
	}
	for worker.Address == "" {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("NetAddress has been updated")
	return nil
}

// setUpNewWorkerCommunication communicates the next node information to the new stage, and the stage before it
func (schedule *Schedule) setUpNewWorkerCommunication(newWorker *types.Worker) {
	if newWorker.Stage != schedule.StageList.MaxPosition {
		for _, worker := range schedule.StageList.FindByPosition(newWorker.Stage + 1).Workers {
			sendNextWorkerAddress(newWorker, worker)
		}
	}
	if newWorker.Stage != 0 {
		for _, worker := range schedule.StageList.FindByPosition(newWorker.Stage - 1).Workers {
			sendNextWorkerAddress(worker, newWorker)
		}
	}
	if newWorker.Stage == 0 {
		message := new(types.Message)
		message.Sender = "0"
		message.Description = common.MsgStartWorker
		connection, err := net.Dial("tcp", newWorker.Address)
		if err != nil {
			panic(err)
		}
		encoder := gob.NewEncoder(connection)
		encoder.Encode(message)
		fmt.Println("Started worker:", newWorker.ID)
	}
}
