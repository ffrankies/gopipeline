package types

import (
	"strconv"
	"sync"
)

// PipelineStageList contains the list of pipeline stages
type PipelineStageList struct {
	List         []*PipelineStage // The list of pipeline stages
	counter      int              // Used to set a counter-type ID to new stages
	counterMutex sync.Mutex       // Ensures that each stage's ID is unique
	MaxPosition  int              // The maximum position of the given stages
}

// NewPipelineStageList creates a new pipeline stage list with an empty list of stages and a nextID of 0
func NewPipelineStageList(numStages int) *PipelineStageList {
	pipelineStageList := new(PipelineStageList)
	pipelineStageList.counter = 0
	for position := 0; position < numStages; position++ {
		stage := newPipelineStage(position)
		pipelineStageList.List = append(pipelineStageList.List, stage)
	}
	pipelineStageList.MaxPosition = numStages - 1
	return pipelineStageList
}

// AddWorker registers a new Worker process with a given PipelineStage
func (stageList *PipelineStageList) AddWorker(host string, position int) *Worker {
	stageList.counterMutex.Lock()
	stageList.counter++
	worker := stageList.FindByPosition(position).AddWorker(strconv.Itoa(stageList.counter), host)
	stageList.counterMutex.Unlock()
	return worker
}

// Length returns the number of stages in the PipelineStageList
func (stageList *PipelineStageList) Length() int {
	return len(stageList.List)
}

// FindStageWithWorker finds the stage that has a worker with the given ID
func (stageList *PipelineStageList) FindStageWithWorker(id string) *PipelineStage {
	for _, stage := range stageList.List {
		for _, worker := range stage.Workers {
			if worker.ID == id {
				return stage
			}
		}
	}
	return nil
}

// FindWorker loops through the workers in the stageList and returns the Worker with a matching ID
func (stageList *PipelineStageList) FindWorker(id string) *Worker {
	for _, stage := range stageList.List {
		for _, worker := range stage.Workers {
			if worker.ID == id {
				return worker
			}
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
			for _, worker := range stage.Workers {
				if worker.PID != -2 && (worker.Address == "" || worker.PID == -1) {
					allUpdated = false
				}
			}
		}
	}
}

// FindBottleneck attempts to find the stage position that executes much slower than its neighbors
func (stageList *PipelineStageList) FindBottleneck() (bottleneckPosition int, scaleNumber int) {
	bottleneckPosition = -1
	bottleneckValue := -1.0
	for position := 0; position <= stageList.MaxPosition; position++ {
		currentPositionExecutionTime := stageList.AverageExecutionTime(position)
		nextPositionExecutionTime := float64(0.0)
		previousPositionExecutionTime := float64(0.0)
		if position != stageList.MaxPosition { // Check if slower than next
			nextPositionExecutionTime = stageList.AverageExecutionTime(position + 1)
			if nextPositionExecutionTime > 0.0 && currentPositionExecutionTime > 1.5*nextPositionExecutionTime {
				difference := currentPositionExecutionTime - nextPositionExecutionTime
				if bottleneckValue < difference {
					bottleneckValue = difference
					bottleneckPosition = position
					scaleNumber = int(currentPositionExecutionTime / nextPositionExecutionTime)
					continue
				}
			}
		}
		if position != 0 { // Check if slower than previous
			previousPositionExecutionTime = stageList.AverageExecutionTime(position - 1)
			if previousPositionExecutionTime > 0.0 && currentPositionExecutionTime > 1.5*previousPositionExecutionTime {
				difference := currentPositionExecutionTime - previousPositionExecutionTime
				if bottleneckValue < difference {
					bottleneckValue = difference
					bottleneckPosition = position
					scaleNumber = int(currentPositionExecutionTime / previousPositionExecutionTime)
					continue
				}
			}
		}
	}
	return
}

// AverageExecutionTime calculates the average execution time given a stage's position
func (stageList *PipelineStageList) AverageExecutionTime(position int) float64 {
	stage := stageList.FindByPosition(position)
	return stage.AverageExecutionTime()
}
