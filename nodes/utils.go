package nodes

import (
	"errors"
	"github.com/mv-orchestration/scheduler/labels"
	"math/rand"
	"reflect"
)

func nodeHasSignificantChanges(oldNode *Node, newNode *Node) bool {
	return oldNode.Name != newNode.Name ||
		oldNode.Labels[labels.NodeCity] != newNode.Labels[labels.NodeCity] ||
		oldNode.Labels[labels.NodeCountry] != newNode.Labels[labels.NodeCountry] ||
		oldNode.Labels[labels.NodeContinent] != newNode.Labels[labels.NodeContinent]
}

func nodeHasAnyLabel(node *Node) bool {
	nodeLabels := [4]string{labels.Node, labels.NodeCity, labels.NodeCountry, labels.NodeContinent}

	for _, label := range nodeLabels {
		if _, ok := node.Labels[label]; ok {
			return true
		}
	}

	return false
}

// GetRandomFromList returns a random node from the list
func GetRandomFromList(options []*Node) (*Node, error) {
	if len(options) == 0 {
		return nil, errors.New("no nodes available")
	}

	return options[rand.Intn(len(options))], nil
}

// GetRandomFromMap returns a random node from the map
func GetRandomFromMap(options map[string][]*Node) (*Node, error) {
	if len(options) == 0 {
		return nil, errors.New("no nodes available")
	}

	keys := reflect.ValueOf(options).MapKeys()
	mapKeyValue := keys[rand.Intn(len(keys))].String()

	return options[mapKeyValue][rand.Intn(len(options[mapKeyValue]))], nil
}
