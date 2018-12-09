package types

import "math"

// PipelineNode struct refers to a computational PipelineNode. A PipelineNode can be assigned multiple functions, or
// pipeline stages.
type PipelineNode struct {
	Address  string    // The internet address of the PipelineNode. Can be DNS, IPv4 or IPv6
	Position int       // The position of this PipelineNode in the PipelineNodelist
	Workers  []*Worker // The workers executing pipeline stages running on this PipelineNode
}

// NewPipelineNode creates a new PipelineNode object
func NewPipelineNode(address string, position int) *PipelineNode {
	pipelineNode := new(PipelineNode)
	pipelineNode.Address = address
	pipelineNode.Position = position
	return pipelineNode
}

// AddWorker adds a PipelineStage to the PipelineNode's workers list
func (pipelineNode *PipelineNode) AddWorker(worker *Worker) {
	pipelineNode.Workers = append(pipelineNode.Workers, worker)
}

// AvailableMemory finds the available memory on this node by finding the minimum AvailableMemory parameter on its
// workers
func (pipelineNode *PipelineNode) AvailableMemory() uint64 {
	minAvailableMemory := uint64(math.MaxUint64)
	for _, worker := range pipelineNode.Workers {
		availableMemory := worker.Stats.NodeAvailableMemory
		if availableMemory < minAvailableMemory {
			minAvailableMemory = availableMemory
		}
	}
	return minAvailableMemory
}

// HasEnoughMemory returns true if the node has enough available memory to contain a worker with the given memory
// requirements
func (pipelineNode *PipelineNode) HasEnoughMemory(requirement uint64) bool {
	availableMemory := pipelineNode.AvailableMemory()
	if availableMemory > requirement {
		return true
	}
	return false
}

// RemoveWorker removes the worker from the Workers list
func (pipelineNode *PipelineNode) RemoveWorker(workerID string) {
	var indexToRemove int
	for index, worker := range pipelineNode.Workers {
		if worker.ID == workerID {
			indexToRemove = index
		}
	}
	if indexToRemove == len(pipelineNode.Workers) {
		pipelineNode.Workers = pipelineNode.Workers[:indexToRemove]
	} else {
		pipelineNode.Workers = append(pipelineNode.Workers[:indexToRemove], pipelineNode.Workers[indexToRemove+1:]...)
	}
}
