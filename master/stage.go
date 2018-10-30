package master

import (
	"strconv"
	"strings"
)

// PipelineStage struct refers to a stage in the pipeline
type PipelineStage struct {
	Host         string // The host on which this stage is being run
	NetAddress   string // The net address to which to
	ListenerPort string // The port number on which the listener is running
	Position     int    // The Stage's position in the pipeline
}

// NewPipelineStage creates a new PipelineStage object. On creation, we don't know the stage's NetAddress or Port, so
// those are initialized as empty strings
func NewPipelineStage(host string, position int) *PipelineStage {
	pipelineStage := new(PipelineStage)
	pipelineStage.Host = host
	pipelineStage.ListenerPort = ""
	pipelineStage.NetAddress = ""
	pipelineStage.Position = position
	return pipelineStage
}

// String converts the PipelineStage struct into a String
func (stage *PipelineStage) String() string {
	pipelineStageString := "PipelineStage {\n"
	pipelineStageString += "\tHost: " + stage.Host + "\n"
	pipelineStageString += "\tNetAddress: " + stage.NetAddress + "\n"
	pipelineStageString += "\tListenerPort: " + stage.ListenerPort + "\n"
	pipelineStageString += "\tPosition: " + strconv.Itoa(stage.Position) + "\n}"
	return pipelineStageString
}

// UpdateListenerPort updates the stage's listener port, and combines the port number and host to form a netaddress
func (stage *PipelineStage) UpdateListenerPort(listenerPort string) {
	stage.ListenerPort = listenerPort
	if strings.Count(stage.Host, ":") > 0 { // If Host is an IPv6 address
		stage.NetAddress = "[" + stage.Host + "]:" + listenerPort
	} else {
		stage.NetAddress = stage.Host + ":" + listenerPort
	}
}
