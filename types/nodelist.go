package types

// PipelineNodeList contains the list of user-provided nodes
type PipelineNodeList struct {
	list []*PipelineNode // The list of nodes
}

// NewPipelineNodeList creates a new empty PipelineNodeList
func NewPipelineNodeList() *PipelineNodeList {
	pipelineNodeList := new(PipelineNodeList)
	return pipelineNodeList
}

// FindNode finds a particular PipelineNode in the PipelineNodeList. If the Node is not found,
// findNode creates a new PipelineNode
func (nodeList *PipelineNodeList) FindNode(nodeAddress string) (pipelineNode *PipelineNode, foundInList bool) {
	for _, node := range nodeList.list {
		if node.Address == nodeAddress {
			pipelineNode = node
			foundInList = true
			return
		}
	}
	pipelineNode = NewPipelineNode(nodeAddress, len(nodeList.list))
	foundInList = false
	return
}

// FindNodeWithEnoughMemory finds a node with enough memory to satisfy the requirement of the given stage
func (nodeList *PipelineNodeList) FindNodeWithEnoughMemory(requirement uint64) *PipelineNode {
	for _, node := range nodeList.list {
		if node.HasEnoughMemory(requirement) {
			return node
		}
	}
	return nil
}

// AddNode will add a PipelineNode to the PipelineNodeList
func (nodeList *PipelineNodeList) AddNode(node *PipelineNode) {
	node.Position = len(nodeList.list) + 1
	nodeList.list = append(nodeList.list, node)
}

// Length returns the number of nodes in the PipelineNodeList
func (nodeList *PipelineNodeList) Length() int {
	return len(nodeList.list)
}

// Pop removes and returns the first PipelineNode from the PipelineNodeList list
func (nodeList *PipelineNodeList) Pop() *PipelineNode {
	firstNode := nodeList.list[0]
	if len(nodeList.list) > 1 {
		nodeList.list = nodeList.list[1:]
	} else {
		nodeList.list = make([]*PipelineNode, 0)
	}
	return firstNode
}
