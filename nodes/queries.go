package nodes

import "github.com/mv-orchestration/scheduler/labels"

// CountNodes returns the number of cluster edge nodes
func (n *Nodes) CountNodes() int {
	return len(n.Nodes)
}

// GetAllNodes list all cluster edge nodes
func (n *Nodes) GetAllNodes() []*Node {
	return n.Nodes
}

// GetNodes list all cluster edge nodes matching filter
func (n *Nodes) GetNodes(filter *NodeFilter) []*Node {
	return n.filterNodes(filter)
}

// AddNode add a new cluster node
func (n *Nodes) AddNode(node *Node) {
	if _, ok := node.Labels[labels.Node]; !ok {
		// Don't add new node if it doesn't have the edge role
		return
	}

	n.Nodes = append(n.Nodes, node)
	n.addToCities(node)
	n.addToCountries(node)
	n.addToContinents(node)
}

// UpdateNode updates a cluster node
func (n *Nodes) UpdateNode(oldNode *Node, newNode *Node) {
	_, oldHasEdgeLabel := oldNode.Labels[labels.Node]
	_, newHasEdgeLabel := newNode.Labels[labels.Node]

	if !oldHasEdgeLabel && newHasEdgeLabel {
		// If node wasn't an edge node but now it is, create it in cache
		n.AddNode(newNode)
	} else if oldHasEdgeLabel && newHasEdgeLabel {
		// If the node is an edge node and has significant update it in cache
		n.updateNodeData(oldNode, newNode)
	} else if oldHasEdgeLabel && !newHasEdgeLabel {
		// If node was an edge node but now it isn't, remove it from cache
		n.DeleteNode(oldNode)
	}
}

// DeleteNode deletes a cluster node
func (n *Nodes) DeleteNode(node *Node) {
	if _, ok := node.Labels[labels.Node]; !ok {
		// Don't try to remove the node if it doesn't have the edge role
		return
	}

	n.removeNodeFromNodes(node)
	n.removeNodeFromCities(node)
	n.removeNodeFromCountries(node)
	n.removeNodeFromContinents(node)
}
