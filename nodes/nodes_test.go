package nodes

import (
	"github.com/aida-dos/gountries"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestNode(name string, edge bool, city string, country string, continent string) *Node {
	labels := map[string]string{
		cityLabel:      city,
		countryLabel:   country,
		continentLabel: continent,
	}

	if edge {
		labels["node-role.kubernetes.io/edge"] = ""
	}

	return &Node{
		Name: name,
		Labels: labels,
	}
}

func newTestNodes() *Nodes {
	return &Nodes{
		Query:          gountries.New(),
		ContinentsList: gountries.NewContinents(),

		Nodes:      make([]*Node, 0),
		Cities:     make(map[string][]*Node),
		Countries:  make(map[string][]*Node),
		Continents: make(map[string][]*Node),
	}
}

func TestNew(t *testing.T) {
	nodes := New()
	assert.Equal(t, 0, nodes.CountNodes())
}

func TestGetAndCountNodes(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", true, "Braga", "Portugal", "Europe")
	nodes.AddNode(node)

	assert.Equal(t, 1, nodes.CountNodes())
	assert.Equal(t, 1, len(nodes.GetAllNodes()))
}

func TestGetNodesFilter(t *testing.T) {
	nodes := newTestNodes()
	filter := &NodeFilter{Resources: Resources{ CPU: 5000 }}

	node0 := &Node{Name: "Node0", CPU: 2500}
	nodes.Nodes = append(nodes.Nodes, node0)

	assert.Equal(t, 0, len(nodes.GetNodes(filter)))

	node1 := &Node{Name: "Node1", CPU: 10000}
	nodes.Nodes = append(nodes.Nodes, node1)

	assert.Equal(t, 1, len(nodes.GetNodes(filter)))
	assert.Equal(t, "Node1", nodes.GetNodes(filter)[0].Name)
}

func TestAddEdgeNodeError(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", true, "RANDOM_C_123", "RANDOM_C_123", "RANDOM_C_123")

	nodes.AddNode(node)
	assert.Equal(t, 1, nodes.CountNodes())
	assert.Equal(t, 0, len(nodes.Cities))
	assert.Equal(t, 0, len(nodes.Countries))
	assert.Equal(t, 0, len(nodes.Continents))
}

func TestAddNormalNode(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", false, "", "", "")

	nodes.AddNode(node)
	assert.Equal(t, 0, nodes.CountNodes())
}

func TestUpdateNodeCoreData(t *testing.T) {
	nodes := newTestNodes()
	oldNode := newTestNode("Node0", true, "Braga", "Portugal", "Europe")

	nodes.AddNode(oldNode)

	assert.Equal(t, 1, len(nodes.Cities["PT-03"]))
	assert.Equal(t, 0, len(nodes.Cities["PT-13"]))

	newNode := newTestNode("Node0", true, "Porto", "Portugal", "Europe")
	nodes.UpdateNode(oldNode, newNode)

	assert.Equal(t, 0, len(nodes.Cities["PT-03"]))
	assert.Equal(t, 1, len(nodes.Cities["PT-13"]))
}

func TestUpdateNodeResources(t *testing.T) {
	nodes := newTestNodes()
	oldNode := newTestNode("Node0", true, "Braga", "Portugal", "Europe")

	nodes.AddNode(oldNode)

	newNode := newTestNode("Node0", true, "Braga", "Portugal", "Europe")
	newNode.Labels["test_label"] = "test"

	nodes.UpdateNode(oldNode, newNode)

	assert.Equal(t, "test", nodes.Nodes[0].Labels["test_label"])
}

func TestDeleteEdgeNode(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", true, "Braga", "Portugal", "Europe")

	nodes.AddNode(node)
	nodes.DeleteNode(node)
	assert.Equal(t, 0, nodes.CountNodes())
}

func TestDeleteNormalNode(t *testing.T) {
	nodes := newTestNodes()
	addNode := newTestNode("Node0", true, "", "", "")
	nodes.AddNode(addNode)

	deleteNode := newTestNode("Node0", false, "", "", "")
	nodes.DeleteNode(deleteNode)

	assert.Equal(t, 1, nodes.CountNodes())
}

func TestDeleteErrorExitsGracefully(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", true, "RANDOM_C_123", "RANDOM_C_123", "RANDOM_C_123")
	nodes.DeleteNode(node)
}

func TestNodeHasSignificantChanges(t *testing.T) {
	oldNode := newTestNode("Node0", true, "Braga", "Portugal", "Europe")
	newNode := newTestNode("Node0", true, "Porto", "Portugal", "Europe")

	assert.Equal(t, false, nodeHasSignificantChanges(oldNode, oldNode))
	assert.Equal(t, true, nodeHasSignificantChanges(oldNode, newNode))
}

func TestFindNodeByName(t *testing.T) {
	nodes := newTestNodes()
	node := newTestNode("Node0", true, "Braga", "Portugal", "Europe")

	nodes.AddNode(node)

	foundNode, _ := nodes.findNodeByName("Node0")
	assert.Equal(t, "Node0", foundNode.Name)

	_, err := nodes.findNodeByName("Node1")
	assert.Error(t, err)
}

func newNodeFilter(
	nodeLabel string, nodeCPU int64, nodeMemory int64,
	filterLabel string, filterCPU int64, filterMemory int64,
) (*Node, *NodeFilter) {

	node := &Node{
		Labels: map[string]string{
			nodeLabel: "test",
		},
		CPU:    nodeCPU,
		Memory: nodeMemory,
	}

	filter := &NodeFilter{
		Labels: map[string]string{filterLabel: "test"},
		Resources: Resources{
			CPU:    filterCPU,
			Memory: filterMemory,
		},
	}

	return node, filter
}

func TestNodeFilter(t *testing.T) {
	node, filter := newNodeFilter("test", 10, 10, "test", 1, 1)
	assert.True(t, nodeMatchesFilters(node, filter))
}

func TestNodeNoFilter(t *testing.T) {
	node, _ := newNodeFilter("test", 10, 10, "test", 1, 1)
	assert.True(t, nodeMatchesFilters(node, &NodeFilter{}))
}

func TestNodeFilterFailLabel(t *testing.T) {
	node, filter := newNodeFilter("test", 10, 10, "test1", 1, 1)
	assert.False(t, nodeMatchesFilters(node, filter))
}

func TestNodeFilterFailCPU(t *testing.T) {
	node, filter := newNodeFilter("test", 0, 10, "test", 1, 1)
	assert.False(t, nodeMatchesFilters(node, filter))
}

func TestNodeFilterFailMemory(t *testing.T) {
	node, filter := newNodeFilter("test", 10, 0, "test", 1, 1)
	assert.False(t, nodeMatchesFilters(node, filter))
}

func TestGetRandomFromListHit(t *testing.T) {
	node := &Node{ Name: "Node0" }
	value, _ := GetRandomFromList([]*Node{node})
	assert.Equal(t, "Node0", value.Name)
}

func TestGetRandomFromListEmptyError(t *testing.T) {
	_, err := GetRandomFromList([]*Node{})
	assert.Error(t, err)
}

func TestGetRandomFromMapHit(t *testing.T) {
	node := &Node{ Name: "Node0" }
	value, _ := GetRandomFromMap(map[string][]*Node{"node": {node}})
	assert.Equal(t, "Node0", value.Name)
}

func TestGetRandomFromMapError(t *testing.T) {
	_, err := GetRandomFromMap(map[string][]*Node{})
	assert.Error(t, err)
}
