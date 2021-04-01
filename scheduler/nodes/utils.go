package nodes

import (
	v1 "k8s.io/api/core/v1"
	"math/rand"
)

func nodeHasSignificantChanges(oldNode *v1.Node, newNode *v1.Node) bool {
	return oldNode.Name != newNode.Name ||
		oldNode.Labels[cityLabel] != newNode.Labels[cityLabel] ||
		oldNode.Labels[countryLabel] != newNode.Labels[countryLabel] ||
		oldNode.Labels[continentLabel] != newNode.Labels[continentLabel]
}

// GetRandom returns a random node from the list
func GetRandom(options []*Node) *Node {
	return options[rand.Intn(len(options))]
}
