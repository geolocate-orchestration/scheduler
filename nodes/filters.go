package nodes

func (n *Nodes) filterNodes(filter *NodeFilter) []*Node {
	filtered := make([]*Node, 0)

	var nodeList []*Node

	if filter.Locations.Cities == nil && filter.Locations.Countries == nil && filter.Locations.Continents == nil {
		nodeList = n.Nodes
	} else {
		nodeList = n.buildFromLocations(filter.Locations)
	}

	for _, n := range nodeList {
		if nodeMatchesFilters(n, filter) {
			filtered = append(filtered, n)
		}
	}
	return filtered
}

func (n *Nodes) buildFromLocations(locations Locations) []*Node {
	nodesMap := make(map[string]*Node)
	n.getCities(locations.Cities, &nodesMap)
	n.getCountries(locations.Countries, &nodesMap)
	n.getContinents(locations.Continents, &nodesMap)

	nodes := make([]*Node, 0, len(nodesMap))
	for  _, value := range nodesMap {
		nodes = append(nodes, value)
	}
	return nodes
}

func (n *Nodes) getCities(locations []string, nodesMap *map[string]*Node) {
	for _, location := range locations {
		for _, node := range n.Cities[location] {
			(*nodesMap)[node.Name] = node
		}
	}
}

func (n *Nodes) getCountries(locations []string, nodesMap *map[string]*Node) {
	for _, location := range locations {
		for _, node := range n.Countries[location] {
			(*nodesMap)[node.Name] = node
		}
	}
}

func (n *Nodes) getContinents(locations []string, nodesMap *map[string]*Node) {
	for _, location := range locations {
		for _, node := range n.Continents[location] {
			(*nodesMap)[node.Name] = node
		}
	}
}

func nodeMatchesFilters(node *Node, filter *NodeFilter) bool {
	if filter == nil {
		return true
	}

	if !nodeHasLabels(node, filter.Labels) {
		return false
	}

	if !matchesResources(node, filter.Resources) {
		return false
	}

	return true
}

func nodeHasLabels(node *Node, labels map[string]string) bool {
	for _, label := range labels {
		if _, ok := node.Labels[label]; !ok || labels[label] != node.Labels[label] {
			return false
		}
	}

	return true
}

func matchesResources(node *Node, resources Resources) bool {
	if resources.CPU != 0 && !nodeHasAllocatableCPU(node, resources.CPU) {
		return false
	}

	if resources.Memory != 0 && !nodeHasAllocatableMemory(node, resources.Memory) {
		return false
	}

	return true
}

func nodeHasAllocatableCPU(node *Node, cpu int64) bool {
	return node.CPU >= cpu
}

func nodeHasAllocatableMemory(node *Node, memory int64) bool {
	return node.Memory >= memory
}
