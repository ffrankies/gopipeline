package types

import (
	"strconv"
)

// PipelineStage struct refers to a stage in the pipeline
type PipelineStage struct {
	Position int       // The Stage's position in the pipeline
	Workers  []*Worker // The list of workers executing this stage
	Scaled   bool      // Marks whether or not this stage has been scaled up or not
}

// NewPipelineStage creates a new PipelineStage object. On creation, we don't know the stage's NetAddress or Port, so
// those are initialized as empty strings
func newPipelineStage(position int) *PipelineStage {
	pipelineStage := new(PipelineStage)
	pipelineStage.Position = position
	pipelineStage.Workers = make([]*Worker, 0)
	pipelineStage.Scaled = false
	return pipelineStage
}

// AddWorker registers a new worker with this pipeline stage
func (stage *PipelineStage) AddWorker(id string, host string) *Worker {
	worker := NewWorker(id, host, stage.Position)
	stage.Workers = append(stage.Workers, worker)
	return worker
}

// AverageExecutionTime calculates the average execution time for the workers in this stage
func (stage *PipelineStage) AverageExecutionTime() float64 {
	totalDuration := float64(0.0)
	for _, worker := range stage.Workers {
		totalDuration += float64(worker.Stats.ExecutionTime)
	}
	return totalDuration / float64(len(stage.Workers))
}

// String converts the PipelineStage struct into a String
func (stage *PipelineStage) String() string {
	pipelineStageString := "PipelineStage {\n"
	pipelineStageString += "\tPosition: " + strconv.Itoa(stage.Position) + "\n}"
	pipelineStageString += "\tNumber of workers running: " + strconv.Itoa(len(stage.Workers)) + "\n"
	pipelineStageString += "\tStage has been scaled: " + strconv.FormatBool(stage.Scaled) + "\n}"
	return pipelineStageString
}
