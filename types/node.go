package types

// PipelineNode struct refers to a computational PipelineNode. A PipelineNode can be assigned multiple functions, or pipeline stages.
type PipelineNode struct {
	Address        string           // The internet address of the PipelineNode. Can be DNS, IPv4 or IPv6
	Position       int              // The position of this PipelineNode in the PipelineNodelist
	PipelineStages []*PipelineStage // The pipeline stages running on this PipelineNode
}

// NewPipelineNode creates a new PipelineNode object
func NewPipelineNode(address string, position int) *PipelineNode {
	pipelineNode := new(PipelineNode)
	pipelineNode.Address = address
	pipelineNode.Position = position
	return pipelineNode
}

// AddStage adds a PipelineStage to the PipelineNode's PipelineStages list
func (pipelineNode *PipelineNode) AddStage(pipelineStage *PipelineStage) {
	pipelineNode.PipelineStages = append(pipelineNode.PipelineStages, pipelineStage)
}
