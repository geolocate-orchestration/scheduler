package geographiclocation

import (
	"aida-scheduler/scheduler/nodes"
	"github.com/aida-dos/gountries"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func newTestGeo(nodes *nodes.Nodes, pod *v1.Pod) *geographiclocation {
	if nodes == nil {
		nodes = newTestNodes(nil, nil, nil, nil)
	}

	return &geographiclocation{
		nodes:      nodes,
		pod:        pod,
		queryType:  "",
		cities:     make([]string, 0),
		countries:  make([]string, 0),
		continents: make([]string, 0),
	}
}

func newTestNodes(
	nodesList []*nodes.Node, citiesList map[string][]*nodes.Node,
	countriesList map[string][]*nodes.Node, continentList map[string][]*nodes.Node,
) *nodes.Nodes {

	if nodesList == nil {
		nodesList = make([]*nodes.Node, 0)
	}

	if citiesList == nil {
		citiesList = make(map[string][]*nodes.Node)
	}

	if countriesList == nil {
		countriesList = make(map[string][]*nodes.Node)
	}

	if continentList == nil {
		continentList = make(map[string][]*nodes.Node)
	}

	return &nodes.Nodes{
		ClientSet: nil,

		Query:          gountries.New(),
		ContinentsList: gountries.NewContinents(),

		Nodes:      nodesList,
		Cities:     citiesList,
		Countries:  countriesList,
		Continents: continentList,
	}
}

func newTestPod(typeString string, value string) *v1.Pod {
	labels := map[string]string{}

	if typeString != "nil" {
		labels["deployment.edge.aida.io/"+typeString+"Location"] = value
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
	}
}

func newTestNode(name string) *nodes.Node {
	return &nodes.Node{
		Name: name,
	}
}

func TestNew(t *testing.T) {
	geoStruct := New(newTestNodes(nil, nil, nil, nil))
	_, err := geoStruct.GetNode(nil)
	assert.Error(t, err, "no nodes are available")
}

func TestGetNodeRandom(t *testing.T) {
	pod := newTestPod("nil", "1")
	nodeStruct := newTestNodes(
		[]*nodes.Node{newTestNode("Node0")},
		nil, nil, nil,
	)
	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodeRequiredFail(t *testing.T) {
	pod := newTestPod("required", "Braga-PT-Europe")
	nodeStruct := newTestNodes(
		[]*nodes.Node{newTestNode("Node0")},
		nil, nil, nil,
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	_, err := geoStruct.GetNode(pod)
	assert.Error(t, err)
}

func TestGetNodeRequiredCityHit(t *testing.T) {
	pod := newTestPod("required", "Braga-PT-Europe")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList,
		map[string][]*nodes.Node{"PT-03": nodeList},
		nil, nil,
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodeRequiredCountryHit(t *testing.T) {
	pod := newTestPod("required", "Braga-PT-Europe")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil,
		map[string][]*nodes.Node{"PT": nodeList},
		nil,
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodeRequiredContinentHit(t *testing.T) {
	pod := newTestPod("required", "Braga-PT-Europe")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil, nil,
		map[string][]*nodes.Node{"EU": nodeList},
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodeRequiredCityDontHitContinent(t *testing.T) {
	pod := newTestPod("required", "Braga-PT-")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil, nil,
		map[string][]*nodes.Node{"EU": nodeList},
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	_, err := geoStruct.GetNode(pod)
	assert.Error(t, err)
}

func TestGetNodePreferredPreferredCityHitCountry(t *testing.T) {
	pod := newTestPod("preferred", "Braga--")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil,
		map[string][]*nodes.Node{"PT": nodeList},
		nil,
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodePreferredPreferredCityHitContinent(t *testing.T) {
	pod := newTestPod("preferred", "Braga--")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil, nil,
		map[string][]*nodes.Node{"EU": nodeList},
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodePreferredCountryHitContinent(t *testing.T) {
	pod := newTestPod("preferred", "-PT-")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil, nil,
		map[string][]*nodes.Node{"EU": nodeList},
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func TestGetNodePreferredNoHit(t *testing.T) {
	pod := newTestPod("preferred", "--")
	nodeList := []*nodes.Node{newTestNode("Node0")}
	nodeStruct := newTestNodes(
		nodeList, nil, nil,
		map[string][]*nodes.Node{"EU": nodeList},
	)

	geoStruct := newTestGeo(nodeStruct, pod)

	node, _ := geoStruct.GetNode(pod)
	assert.Equal(t, "Node0", node.Name)
}

func getLocationLabel(t *testing.T, typeString string) {
	pod := newTestPod(typeString, "1")
	geoStruct := newTestGeo(nil, pod)
	r := geoStruct.getLocationLabelType()
	assert.Equal(t, typeString, r)
}

func TestGetLocationLabelRequired(t *testing.T) {
	getLocationLabel(t, "required")
}

func TestGetLocationLabelPreferred(t *testing.T) {
	getLocationLabel(t, "preferred")
}

func TestGetLocationLabelNone(t *testing.T) {
	getLocationLabel(t, "")
}

func TestGetRandomNodeEmpty(t *testing.T) {
	_, err := getRandomNode(newTestNodes(nil, nil, nil, nil))
	assert.Error(t, err, "no nodes are available")
}

func TestGetRandomNode(t *testing.T) {
	nodesStruct := newTestNodes(nil, nil, nil, nil)
	nodesStruct.Nodes = append(nodesStruct.Nodes, newTestNode("Node0"))
	node, _ := getRandomNode(nodesStruct)
	assert.Equal(t, "Node0", node.Name)
}

func TestParseLocations(t *testing.T) {
	geoStruct := newTestGeo(nil, nil)
	locationLabel := "Braga_Porto_Madrid-PT-Europe"
	geoStruct.parseLocations(locationLabel)

	assert.Equal(t, 3, len(geoStruct.cities))
	assert.Equal(t, 1, len(geoStruct.countries))
	assert.Equal(t, 1, len(geoStruct.continents))
}
