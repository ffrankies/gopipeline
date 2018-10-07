package master

import (
	"strconv"
	"strings"
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
	var stringBuilder strings.Builder
	stringBuilder.WriteString("PipelineStage {\n")
	stringBuilder.WriteString("\tNodeAddress: " + stage.NodeAddress + "\n")
	stringBuilder.WriteString("\tMasterSocketPort: " + strconv.Itoa(stage.MasterSocketPort) + "\n")
	stringBuilder.WriteString("\tInputSocketPort: " + strconv.Itoa(stage.InputSocketPort) + "\n")
	stringBuilder.WriteString("\tOutputSocketPort: " + strconv.Itoa(stage.OutputSocketPort) + "\n")
	stringBuilder.WriteString("\tPosition: " + strconv.Itoa(stage.Position) + "\n")
	stringBuilder.WriteString("}")
	return stringBuilder.String()
}
