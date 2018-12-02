package scheduler

import (
	"fmt"
	"math"
	"time"

	"github.com/ffrankies/gopipeline/types"
)

// Schedule contains the information needed for scheduling
type Schedule struct {
	freeNodeList *types.PipelineNodeList  // The list of Nodes available for scheduling
	NodeList     *types.PipelineNodeList  // The list of Nodes that have at least one stages running on them
	StageList    *types.PipelineStageList // The list of pipeline Stages, with metadata
}

// NewSchedule creates a new scheduler with empty node and stage lists, and populates the empty node list
func NewSchedule(nodeList []string) *Schedule {
	schedule := new(Schedule)
	schedule.NodeList = types.NewPipelineNodeList()
	schedule.StageList = types.NewPipelineStageList()
	schedule.freeNodeList = types.NewPipelineNodeList()
	for _, nodeHostName := range nodeList {
		node := types.NewPipelineNode(nodeHostName, -1)
		schedule.freeNodeList.AddNode(node)
	}
	return schedule
}

// Static does initial static scheduling of the pipeline stages on the available nodes
func (schedule *Schedule) Static(functionList []types.AnyFunc) {
	fmt.Println("Performing static scheduling")
	density := schedule.CalculateFunctionDensity(functionList)
	counter := 0
	schedulingNode := schedule.freeNodeList.Pop()
	for _, function := range functionList {
		schedule.AssignStageToNode(function, schedulingNode)
		counter++
		if counter == density {
			counter = 0
			density = schedule.CalculateFunctionDensity(functionList)
			if density < 0 { // If scheduling is over, there are no functions to schedule, so density becomes < 0
				break
			}
			schedulingNode = schedule.freeNodeList.Pop()
		}
	}
}

// CalculateFunctionDensity calculates the initial function density in the pipeline
func (schedule *Schedule) CalculateFunctionDensity(functionList []types.AnyFunc) int {
	numFunctionsRemaining := len(functionList) - schedule.StageList.Length()
	numNodesRemaining := schedule.freeNodeList.Length()
	density := math.Ceil(float64(numFunctionsRemaining) / float64(numNodesRemaining))
	return int(density)
}

// AssignStageToFreeNode assigns a single pipeline stage to a single free node
func (schedule *Schedule) AssignStageToFreeNode(function types.AnyFunc) {
	if schedule.freeNodeList.Length() == 0 {
		panic("FATAL ERROR: There are no free nodes to assign this stage to")
	}
	schedulingNode := schedule.freeNodeList.Pop()
	schedule.AssignStageToNode(function, schedulingNode)
}

// AssignStageToNode assigns a single pipeline stage (function) to a single node
func (schedule *Schedule) AssignStageToNode(function types.AnyFunc, pipelineNode *types.PipelineNode) {
	_, foundInList := schedule.NodeList.FindNode(pipelineNode.Address)
	pipelineStage := schedule.StageList.AddStage(pipelineNode.Address, schedule.StageList.Length())
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		schedule.NodeList.AddNode(pipelineNode)
	}
}

// UpdateStageStats updates the worker statistics for a given stage from an incoming message
func (schedule *Schedule) UpdateStageStats(message *types.Message) {
	fmt.Println("Received worker stats from", message.Sender)
	stage := schedule.StageList.Find(message.Sender)
	stageStats, ok := (message.Contents).(*types.WorkerStats)
	if ok {
		stage.Stats = stageStats
	} else {
		fmt.Println("ERROR: Could not convert message contents to WorkerStats")
	}
}

// UpdateStageInfo updates the stage information for a given stage from an incoming message
func (schedule *Schedule) UpdateStageInfo(message *types.Message) {
	fmt.Println("Received worker info from", message.Sender)
	stage := schedule.StageList.Find(message.Sender)
	stageInfo, ok := (message.Contents).(types.MessageStageInfo)
	if ok {
		stage.NetAddress = stageInfo.Address
		stage.PID = stageInfo.PID
	} else {
		fmt.Println("ERROR: Could not convert message contents to MessageStageInfo")
	}
}

// Dynamic does dynamic scheduling of the pipeline stages on the available nodes, with the aim of increasing
// throughput and memory utilization
func (schedule *Schedule) Dynamic() {
	for {
		time.Sleep(1 * time.Second)
	}
}
