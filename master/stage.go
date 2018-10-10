package master

import (
	"strconv"
)

// PipelineStage struct refers to a stage in the pipeline
type PipelineStage struct {
	NodeAddress      string // The IP address of the node on which this is run
	MasterSocketPort int    // The port number for the worker socket on this node
	InputSocketPort  int    // The port number for the input socket on this node
	OutputSocketPort int    // The port number for the output socket on this node
	Position         int    // The Stage's position in the pipeline
}

// NewPipelineStage creates a new PipelineStage object
func NewPipelineStage(nodeAddress string, masterSocketPort int, inputSocketPort int, outputSocketPort int,
	position int) *PipelineStage {
	pipelineStage := new(PipelineStage)
	pipelineStage.NodeAddress = nodeAddress
	pipelineStage.MasterSocketPort = masterSocketPort
	pipelineStage.InputSocketPort = inputSocketPort
	pipelineStage.OutputSocketPort = outputSocketPort
	pipelineStage.Position = position
	return pipelineStage
}

// String converts the PipelineStage struct into a String
func (stage *PipelineStage) String() string {
	pipelineStageString := "PipelineStage {\n"
	pipelineStageString += "\tNodeAddress: " + stage.NodeAddress + "\n"
	pipelineStageString += "\tMasterSocketPort: " + strconv.Itoa(stage.MasterSocketPort) + "\n"
	pipelineStageString += "\tInputSocketPort: " + strconv.Itoa(stage.InputSocketPort) + "\n"
	pipelineStageString += "\tOutputSocketPort: " + strconv.Itoa(stage.OutputSocketPort) + "\n"
	pipelineStageString += "\tPosition: " + strconv.Itoa(stage.Position) + "\n}"
	return pipelineStageString
}
