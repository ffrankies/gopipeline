package scheduler

import (
	"math"
	"time"

	"github.com/ffrankies/gopipeline/types"
)

// PipelineStageList is the list of pipeline stages
var PipelineStageList = types.NewPipelineStageList()

// PipelineNodeList is the list of pipeline nodes
var PipelineNodeList []*types.PipelineNode

// Static does initial static scheduling of the pipeline stages on the available nodes
func Static(functionList []types.AnyFunc, nodeList []string) {
	density := calculateFunctionDensity(functionList, nodeList)
	counter := 0
	nodeIndex := 0
	for _, function := range functionList {
		assignStageToNode(function, nodeList[nodeIndex])
		counter++
		if counter == density {
			nodeIndex++
			counter = 0
			density = calculateFunctionDensity(functionList, nodeList)
		}
	}
}

// assignStageToNode assigns a single pipeline stage (function) to a single node
func assignStageToNode(function types.AnyFunc, nodeAddress string) {
	pipelineNode, foundInList := findNode(nodeAddress)
	pipelineStage := PipelineStageList.AddStage(nodeAddress, PipelineStageList.Length())
	pipelineNode.AddStage(pipelineStage)
	if foundInList == false {
		PipelineNodeList = append(PipelineNodeList, pipelineNode)
	}
}

// findNode finds a particular PipelineNode in the PipelineNodeList. If the Node is not found,
// findNode creates a new PipelineNode
func findNode(nodeAddress string) (pipelineNode *types.PipelineNode, foundInList bool) {
	for _, node := range PipelineNodeList {
		if node.Address == nodeAddress {
			pipelineNode = node
			foundInList = true
			return
		}
	}
	pipelineNode = types.NewPipelineNode(nodeAddress, len(PipelineNodeList))
	foundInList = false
	return
}

// calculateFunctionDensity calculates the initial function density in the pipeline
func calculateFunctionDensity(functionList []types.AnyFunc, nodeList []string) int {
	numFunctions := len(functionList) - PipelineStageList.Length()
	numNodes := len(nodeList) - len(PipelineNodeList)
	density := math.Ceil(float64(numFunctions) / float64(numNodes))
	return int(density)
}

// Dynamic does dynamic scheduling of the pipeline stages on the available nodes, with the aim of increasing
// throughput and memory utilization
func Dynamic() {
	for {
		time.Sleep(1 * time.Second)
	}
}
