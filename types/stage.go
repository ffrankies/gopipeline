package types

import (
	"strconv"
)

// PipelineStage struct refers to a stage in the pipeline
type PipelineStage struct {
	Host       string       // The host on which this stage is being run
	NetAddress string       // The net address to which to
	Position   int          // The Stage's position in the pipeline
	StageID    string       // The ID of this stage
	PID        int          // The PID of the worker process running this stage
	Stats      *WorkerStats // The performance statistics of the worker process running this stage
}

// NewPipelineStage creates a new PipelineStage object. On creation, we don't know the stage's NetAddress or Port, so
// those are initialized as empty strings
func newPipelineStage(host string, position int, stageID string) *PipelineStage {
	pipelineStage := new(PipelineStage)
	pipelineStage.Host = host
	pipelineStage.NetAddress = ""
	pipelineStage.PID = -1
	pipelineStage.Position = position
	pipelineStage.StageID = stageID
	return pipelineStage
}

// String converts the PipelineStage struct into a String
func (stage *PipelineStage) String() string {
	pipelineStageString := "PipelineStage {\n"
	pipelineStageString += "\tHost: " + stage.Host + "\n"
	pipelineStageString += "\tNetAddress: " + stage.NetAddress + "\n"
	pipelineStageString += "\tPID: " + strconv.Itoa(stage.PID) + "\n"
	pipelineStageString += "\tPosition: " + strconv.Itoa(stage.Position) + "\n}"
	pipelineStageString += "\tStageID: " + stage.StageID + "\n"
	pipelineStageString += "\tStats: " + stage.Stats.String() + "\n}"
	return pipelineStageString
}
