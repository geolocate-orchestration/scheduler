package nodes

import (
	"errors"
	"fmt"
	"github.com/geolocate-orchestration/gountries"
	"github.com/geolocate-orchestration/scheduler/labels"
	"k8s.io/klog/v2"
)

func (n *Nodes) addToCities(node *Node) {
	cityValue := node.Labels[labels.NodeCity]

	if cityValue != "" {
		if city, err := n.Query.FindSubdivisionByName(cityValue); err == nil {
			cityCode := fmt.Sprintf("%s-%s", city.CountryAlpha2, city.Code)
			n.Cities[cityCode] = append(n.Cities[cityCode], node)
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) addToCountries(node *Node) {
	countryValue := node.Labels[labels.NodeCountry]

	if countryValue != "" {
		if country, err := n.findCountry(countryValue); err == nil {
			n.Countries[country.Alpha2] = append(n.Countries[country.Alpha2], node)
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) addToContinents(node *Node) {
	continentValue := node.Labels[labels.NodeContinent]

	if continentValue != "" {
		if continent, err := n.ContinentsList.FindContinent(continentValue); err == nil {
			n.Continents[continent.Code] = append(n.Continents[continent.Code], node)
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) updateNodeData(savedNode *Node, newNode *Node) {
	if nodeHasSignificantChanges(savedNode, newNode) {
		n.DeleteNode(savedNode)
		n.AddNode(newNode)
		klog.Infof("node replaced in cache: %s\n", savedNode.Name)
	} else {
		n.updateNodeFields(savedNode, newNode)
	}
}

func (n *Nodes) updateNodeFields(savedNode *Node, newNode *Node) {
	if savedNode.CPU != newNode.CPU {
		klog.Infof("updated node %s CPU: %d -> %d\n", savedNode.Name, savedNode.CPU, newNode.CPU)
		savedNode.CPU = newNode.CPU
	}

	if savedNode.Memory != newNode.Memory {
		klog.Infof("updated node %s Memory: %d -> %d\n", savedNode.Name, savedNode.Memory, newNode.Memory)
		savedNode.Memory = newNode.Memory
	}

	for key, newValue := range newNode.Labels {
		oldValue, ok := savedNode.Labels[key]

		if !ok {
			savedNode.Labels[key] = newValue
			klog.Infof("updated node %s - added label '%s' with value '%s'\n", savedNode.Name, key, newValue)
		} else if oldValue != newValue {
			savedNode.Labels[key] = newValue
			klog.Infof("updated node %s - '%s' label changed: '%s' -> '%s'\n", savedNode.Name, key, oldValue, newValue)
		}
	}

	for key := range savedNode.Labels {
		_, ok := newNode.Labels[key]

		if !ok {
			delete(savedNode.Labels, key)
			klog.Infof("updated node %s - deleted label '%s'\n", savedNode.Name, key)
		}
	}
}

func (n *Nodes) findNodeByName(name string) (*Node, error) {
	for _, node := range n.Nodes {
		if node.Name == name {
			return node, nil
		}
	}

	return nil, errors.New("node with given name not found")
}

func (n *Nodes) removeNodeFromNodes(node *Node) {
	for i, v := range n.Nodes {
		if v.Name == node.Name {
			n.Nodes = append(n.Nodes[:i], n.Nodes[i+1:]...)
		}
	}
}

func (n *Nodes) removeNodeFromCities(node *Node) {
	cityValue := node.Labels[labels.NodeCity]

	if cityValue != "" {
		if city, err := n.Query.FindSubdivisionByName(cityValue); err == nil {
			cityCode := fmt.Sprintf("%s-%s", city.CountryAlpha2, city.Code)
			for i, v := range n.Cities[cityCode] {
				if v.Name == node.Name {
					n.Cities[cityCode] = append(n.Cities[cityCode][:i], n.Cities[cityCode][i+1:]...)
				}
			}
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) removeNodeFromCountries(node *Node) {
	countryValue := node.Labels[labels.NodeCountry]
	if countryValue != "" {
		if country, err := n.findCountry(countryValue); err == nil {
			for i, v := range n.Countries[country.Alpha2] {
				if v.Name == node.Name {
					n.Countries[country.Alpha2] =
						append(n.Countries[country.Alpha2][:i], n.Countries[country.Alpha2][i+1:]...)
				}
			}
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) removeNodeFromContinents(node *Node) {
	continentValue := node.Labels[labels.NodeContinent]

	if continentValue != "" {
		if continent, err := n.ContinentsList.FindContinent(continentValue); err == nil {
			for i, v := range n.Continents[continent.Code] {
				if v.Name == node.Name {
					n.Continents[continent.Code] =
						append(n.Continents[continent.Code][:i], n.Continents[continent.Code][i+1:]...)
				}
			}
		} else {
			klog.Errorln(err)
		}
	}
}

func (n *Nodes) findCountry(countryID string) (gountries.Country, error) {
	if country, err := n.Query.FindCountryByName(countryID); err == nil {
		return country, nil
	}

	if country, err := n.Query.FindCountryByAlpha(countryID); err == nil {
		return country, nil
	}

	return gountries.Country{}, errors.New("given country identifier does not match any country")
}
