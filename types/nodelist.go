package types

// PipelineNodeList contains the List of user-provided nodes
type PipelineNodeList struct {
	List []*PipelineNode // The List of nodes
}

// NewPipelineNodeList creates a new empty PipelineNodeList
func NewPipelineNodeList() *PipelineNodeList {
	pipelineNodeList := new(PipelineNodeList)
	return pipelineNodeList
}

// FindNode finds a particular PipelineNode in the PipelineNodeList. If the Node is not found,
// findNode creates a new PipelineNode
func (nodeList *PipelineNodeList) FindNode(nodeAddress string) (pipelineNode *PipelineNode, foundInList bool) {
	for _, node := range nodeList.List {
		if node.Address == nodeAddress {
			pipelineNode = node
			foundInList = true
			return
		}
	}
	pipelineNode = NewPipelineNode(nodeAddress, len(nodeList.List))
	foundInList = false
	return
}

// FindNodeWithEnoughMemory finds a node with enough memory to satisfy the requirement of the given stage
func (nodeList *PipelineNodeList) FindNodeWithEnoughMemory(requirement uint64) *PipelineNode {
	for _, node := range nodeList.List {
		if node.HasEnoughMemory(requirement) {
			return node
		}
	}
	return nil
}

// AddNode will add a PipelineNode to the PipelineNodeList
func (nodeList *PipelineNodeList) AddNode(node *PipelineNode) {
	node.Position = len(nodeList.List) + 1
	nodeList.List = append(nodeList.List, node)
}

// Length returns the number of nodes in the PipelineNodeList
func (nodeList *PipelineNodeList) Length() int {
	return len(nodeList.List)
}

// Pop removes and returns the first PipelineNode from the PipelineNodeList List
func (nodeList *PipelineNodeList) Pop() *PipelineNode {
	firstNode := nodeList.List[0]
	if len(nodeList.List) > 1 {
		nodeList.List = nodeList.List[1:]
	} else {
		nodeList.List = make([]*PipelineNode, 0)
	}
	return firstNode
}
