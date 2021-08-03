package nodes

import (
	"k8s.io/klog/v2"
)

// CountNodes returns the number of cluster nodes
func (n *Nodes) CountNodes() int {
	return len(n.Nodes)
}

// GetAllNodes list all cluster nodes
func (n *Nodes) GetAllNodes() []*Node {
	return n.Nodes
}

// GetNodes list all cluster nodes matching filter
func (n *Nodes) GetNodes(filter *NodeFilter) []*Node {
	return n.filterNodes(filter)
}

// AddNode add a new cluster node
func (n *Nodes) AddNode(node *Node) {
	if !nodeHasAnyLabel(node) {
		// Don't add new node if it doesn't have the node.geolocate.io role
		return
	}

	n.Nodes = append(n.Nodes, node)
	n.addToCities(node)
	n.addToCountries(node)
	n.addToContinents(node)
	klog.Infof("node added to cache: %s\n", node.Name)
}

// UpdateNode updates a cluster node
func (n *Nodes) UpdateNode(oldNode *Node, newNode *Node) {
	savedNode, err := n.findNodeByName(oldNode.Name)
	if err != nil {
		// If node wasn't labeled but now it is, create it in cache
		n.AddNode(newNode)
		return
	}

	oldHasNodeLabel := nodeHasAnyLabel(savedNode)
	newHasNodeLabel := nodeHasAnyLabel(newNode)

	if !oldHasNodeLabel && newHasNodeLabel {
		// If node wasn't labeled but now it is, create it in cache
		n.AddNode(newNode)
	} else if oldHasNodeLabel && newHasNodeLabel {
		// If the node is labeled and has significant update it in cache
		n.updateNodeData(savedNode, newNode)
	} else if oldHasNodeLabel && !newHasNodeLabel {
		// If node was labeled but now it isn't, remove it from cache
		n.DeleteNode(savedNode)
	}
}

// DeleteNode deletes a cluster node
func (n *Nodes) DeleteNode(node *Node) {
	n.removeNodeFromNodes(node)
	n.removeNodeFromCities(node)
	n.removeNodeFromCountries(node)
	n.removeNodeFromContinents(node)
	klog.Infof("node deleted from cache: %s\n", node.Name)
}
