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
		newStage := schedule.AssignStageToFreeNode(position)
		schedule.startStage(newStage, program, masterAddress)
		fmt.Println("Waiting for worker to send info...")
		if err := schedule.waitForWorkerToSendInfo(newStage); err != nil {
			panic(err)
		}
		fmt.Println("Done waiting for worker to send info...")
		schedule.setUpNewWorkerCommunication(newStage)
		numScaled++
	}
	// TODO(): also scale on underutilized nodes
}

// waitForWorkerToSendInfo busy waits until the worker sends its info
func (schedule *Schedule) waitForWorkerToSendInfo(stage *types.PipelineStage) error {
	for stage.PID == -1 {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("PID has been updated")
	if stage.PID == -2 {
		return errors.New("ERROR: Stage could not be started")
	}
	for stage.NetAddress == "" {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("NetAddress has been updated")
	return nil
}

// setUpNewWorkerCommunication communicates the next node information to the new stage, and the stage before it
func (schedule *Schedule) setUpNewWorkerCommunication(newStage *types.PipelineStage) {
	if newStage.Position != schedule.StageList.MaxPosition {
		for _, stage := range schedule.StageList.List {
			if stage.Position == newStage.Position+1 {
				sendNextWorkerAddress(newStage, stage)
			}
		}
		for _, stage := range schedule.StageList.List {
			if stage.Position == newStage.Position-1 {
				sendNextWorkerAddress(stage, newStage)
			}
		}
	}
	if newStage.Position == 0 {
		message := new(types.Message)
		message.Sender = "0"
		message.Description = common.MsgStartWorker
		connection, err := net.Dial("tcp", newStage.NetAddress)
		if err != nil {
			panic(err)
		}
		encoder := gob.NewEncoder(connection)
		encoder.Encode(message)
		fmt.Println("Started stage:", newStage.StageID)
	}
}
