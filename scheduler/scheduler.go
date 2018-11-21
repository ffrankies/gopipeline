package scheduler

import (
	"math"
	"time"

	"github.com/ffrankies/gopipeline/types"
)

// Schedule contains the information needed for scheduling
type Schedule struct {
	NodeList  *types.PipelineNodeList  // The list of available nodes, with metadata
	StageList *types.PipelineStageList // The list of pipeline stages, with metadata
}

// NewSchedule creates a new scheduler with empty node and stage lists
func NewSchedule() *Schedule {
	schedule := new(Schedule)
	schedule.NodeList = types.NewPipelineNodeList()
	schedule.StageList = types.NewPipelineStageList()
	return schedule
}

// AssignStageToNode assigns a single pipeline stage (function) to a single node
func (schedule *Schedule) AssignStageToNode(function types.AnyFunc, nodeAddress string) {
	pipelineNode, foundInList := schedule.NodeList.FindNode(nodeAddress)
	pipelineStage := schedule.StageList.AddStage(nodeAddress, schedule.StageList.Length())
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		schedule.NodeList.AddNode(pipelineNode)
	}
}

// Static does initial static scheduling of the pipeline stages on the available nodes
func (schedule *Schedule) Static(functionList []types.AnyFunc, nodeList []string) {
	density := schedule.CalculateFunctionDensity(functionList, nodeList)
	counter := 0
	nodeIndex := 0
	for _, function := range functionList {
		schedule.AssignStageToNode(function, nodeList[nodeIndex])
		counter++
		if counter == density {
			nodeIndex++
			counter = 0
			density = schedule.CalculateFunctionDensity(functionList, nodeList)
		}
	}
}

// CalculateFunctionDensity calculates the initial function density in the pipeline
func (schedule *Schedule) CalculateFunctionDensity(functionList []types.AnyFunc, nodeList []string) int {
	numFunctionsRemaining := len(functionList) - schedule.StageList.Length()
	numNodesRemaining := len(nodeList) - schedule.NodeList.Length()
	density := math.Ceil(float64(numFunctionsRemaining) / float64(numNodesRemaining))
	return int(density)
}

// Dynamic does dynamic scheduling of the pipeline stages on the available nodes, with the aim of increasing
// throughput and memory utilization
func Dynamic() {
	for {
		time.Sleep(1 * time.Second)
	}
}
