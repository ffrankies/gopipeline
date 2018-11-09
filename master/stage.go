package master

import (
	"strconv"
	"sync"
)

// PipelineStageList contains the list of pipeline stages
type PipelineStageList struct {
	List         []*PipelineStage // The list of pipeline stages
	counter      int              // Used to set a counter-type ID to new stages
	counterMutex sync.Mutex       // Ensures that each stage's ID is unique
}

// NewPipelineStageList creates a new pipeline stage list with an empty list of stages and a nextID of 0
func NewPipelineStageList() *PipelineStageList {
	pipelineStageList := new(PipelineStageList)
	pipelineStageList.counter = 0
	return pipelineStageList
}

// AddStage adds a new PipelineStage to the PipelineStageList
func (stageList *PipelineStageList) AddStage(host string, position int) *PipelineStage {
	stageList.counterMutex.Lock()
	stageList.counter++
	stage := newPipelineStage(host, position, strconv.Itoa(stageList.counter))
	stageList.counterMutex.Unlock()
	stageList.List = append(stageList.List, stage)
	return stage
}

// Length returns the number of stages in the PipelineStageList
func (stageList *PipelineStageList) Length() int {
	return len(stageList.List)
}

// Find loops through the stageList and returns the PipelineStage with a matching ID
func (stageList *PipelineStageList) Find(id string) *PipelineStage {
	for _, stage := range stageList.List {
		if stage.StageID == id {
			return stage
		}
	}
	return nil
}

// FindByPosition loops through the stageList and returns the PipelineStage with a matching position. IN the future, it
// should return a list of PipelineStages with a matching position
func (stageList *PipelineStageList) FindByPosition(position int) *PipelineStage {
	for _, stage := range stageList.List {
		if stage.Position == position {
			return stage
		}
	}
	return nil
}

// WaitUntilAllListenerPortsUpdated busy waits until there are no more stages whose listener port needs to be updated
func (stageList *PipelineStageList) WaitUntilAllListenerPortsUpdated() {
	allUpdated := false
	for allUpdated == false {
		allUpdated = true
		for _, stage := range stageList.List {
			if stage.NetAddress == "" || stage.PID == -1 {
				allUpdated = false
			}
		}
	}
}

// PipelineStage struct refers to a stage in the pipeline
type PipelineStage struct {
	Host       string // The host on which this stage is being run
	NetAddress string // The net address to which to
	Position   int    // The Stage's position in the pipeline
	StageID    string // The ID of this stage
	PID        int    // The PID of the worker process running this stage
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
	pipelineStageString += "\tStageID: " + stage.StageID + "\n}"
	return pipelineStageString
}
