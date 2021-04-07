package nodes

import (
	"errors"
	"math/rand"
	"reflect"
)

func nodeHasSignificantChanges(oldNode *Node, newNode *Node) bool {
	return oldNode.Name != newNode.Name ||
		oldNode.Labels[cityLabel] != newNode.Labels[cityLabel] ||
		oldNode.Labels[countryLabel] != newNode.Labels[countryLabel] ||
		oldNode.Labels[continentLabel] != newNode.Labels[continentLabel]
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
