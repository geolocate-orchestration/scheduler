package location

import (
	"errors"
	"fmt"
	"github.com/mv-orchestration/gountries"
	"github.com/mv-orchestration/scheduler/algorithms"
	"github.com/mv-orchestration/scheduler/labels"
	"github.com/mv-orchestration/scheduler/nodes"
	"k8s.io/klog/v2"
	"strings"
)

type location struct {
	query      *gountries.Query
	nodes      nodes.INodes
	pod        *algorithms.Workload
	queryType  string // required or preferred
	cities     []string
	countries  []string
	continents []string
}

// New creates new location struct
func New(nodes nodes.INodes) algorithms.Algorithm {
	return &location{
		query:      gountries.New(),
		nodes:      nodes,
		pod:        nil,
		queryType:  "",
		cities:     make([]string, 0),
		countries:  make([]string, 0),
		continents: make([]string, 0),
	}
}

func (g *location) GetName() string {
	return "location"
}

// GetNode select the best node matching the given constraints labels
// It returns error if there are no nodes available and if no node matches an existing 'requiredLocation' label
func (g *location) GetNode(pod *algorithms.Workload) (*nodes.Node, error) {
	var node *nodes.Node
	var err error

	if g.nodes.CountNodes() == 0 {
		errMessage := "no nodes are available"
		return nil, errors.New(errMessage)
	}

	g.pod = pod
	if queryType := g.getLocationLabelType(); queryType != "" {
		g.queryType = queryType
		node, err = g.getNodeByLocation()
	} else {
		// Node location labels were set so returning a random node
		node, err = nodes.GetRandomFromList(g.nodes.GetNodes(g.buildResourceFilter()))
	}

	return node, err
}

// Locations

func (g *location) buildResourceFilter() *nodes.NodeFilter {
	return &nodes.NodeFilter{
		Labels: nil,
		Resources: nodes.Resources{
			CPU:    g.pod.CPU,
			Memory: g.pod.Memory,
		},
	}
}

func (g *location) getNodeByLocation() (*nodes.Node, error) {
	label := ""
	switch g.queryType {
	case "required":
		label = labels.WorkloadRequiredLocation
	case "preferred":
		label = labels.WorkloadPreferredLocation
	}

	locations := g.pod.Labels[label]
	klog.Infoln(label, locations)

	// fill location info from labels in the geo struct
	g.parseLocations(locations)

	if node, err := g.getRequestedLocation(); err == nil {
		return node, nil
	} else if g.queryType == "required" {
		// if location is "required" but there are no matching nodes, throw error
		return nil, err
	}

	if node, err := g.getSimilarToRequestedLocation(); err == nil {
		return node, nil
	}

	// when location is "preferred" and there are no matching nodes, return random node
	return nodes.GetRandomFromList(g.nodes.GetNodes(g.buildResourceFilter()))
}

func (g *location) getRequestedLocation() (*nodes.Node, error) {
	if node, err := g.getByCity(); err == nil {
		return node, nil
	}

	if node, err := g.getByCountry(); err == nil {
		return node, nil
	}

	if node, err := g.getByContinent(); err == nil {
		return node, nil
	}

	return nil, errors.New("no nodes match given locations")
}

func (g *location) getSimilarToRequestedLocation() (*nodes.Node, error) {
	countries := make(map[string]bool)
	continents := make(map[string]bool)

	g.getCitiesPredecessors(g.cities, &countries, &continents)
	g.getCountriesPredecessors(g.countries, &continents)

	if options := getNodes(g.nodes, g.pod, nil, getKeys(countries), getKeys(continents)); len(options) > 0 {
		return nodes.GetRandomFromList(options)
	}

	return nil, errors.New("no nodes match similar location to given locations")
}

// GetBy

func (g *location) getByCity() (*nodes.Node, error) {
	cities := make([]string, 0)

	for _, city := range g.cities {
		city, err := g.query.FindSubdivisionByName(city)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		cityCode := fmt.Sprintf("%s-%s", city.CountryAlpha2, city.Code)
		cities = append(cities, cityCode)
	}

	options := getNodes(g.nodes, g.pod, cities, nil, nil)
	return nodes.GetRandomFromList(options)
}

func (g *location) getByCountry() (*nodes.Node, error) {
	countries := make([]string, 0)

	for _, countryName := range g.countries {
		country, err := g.findCountry(countryName)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		countries = append(countries, country.Alpha2)
	}

	options := getNodes(g.nodes, g.pod, nil, countries, nil)
	return nodes.GetRandomFromList(options)
}

func (g *location) getByContinent() (*nodes.Node, error) {
	continents := make([]string, 0)
	gcont := gountries.NewContinents()

	for _, continentID := range g.continents {
		if continent, err := gcont.FindContinent(continentID); err == nil {
			continents = append(continents, continent.Code)
		}
	}

	options := getNodes(g.nodes, g.pod, nil, nil, continents)
	return nodes.GetRandomFromList(options)
}

// Helpers

func getNodes(inodes nodes.INodes, workload *algorithms.Workload, cities []string, countries []string, continents []string) []*nodes.Node {
	nodeFilter := &nodes.NodeFilter{
		Locations: nodes.Locations{
			Cities:     cities,
			Countries:  countries,
			Continents: continents,
		},
		Resources: nodes.Resources{
			CPU:    workload.CPU,
			Memory: workload.Memory,
		},
	}

	return inodes.GetNodes(nodeFilter)
}

func (g *location) getCitiesPredecessors(cities []string, countries *map[string]bool, continents *map[string]bool) {
	for _, city := range cities {
		country, err := g.query.FindSubdivisionCountryByName(city)
		if err != nil {
			// If subdivision name does not exists skip
			continue
		}
		if _, ok := (*countries)[country.Alpha2]; ok {
			// If country already processed skip
			continue
		}

		(*countries)[country.Alpha2] = true
		(*continents)[country.Continent] = true
	}
}

func (g *location) getCountriesPredecessors(countries []string, continents *map[string]bool) {
	for _, country := range countries {
		country, err := g.findCountry(country)
		if err != nil {
			// If country name does not exists skip
			continue
		}
		(*continents)[country.Continent] = true
	}
}

func (g *location) getLocationLabelType() string {
	if g.pod.Labels[labels.WorkloadRequiredLocation] != "" {
		return "required"
	}

	if g.pod.Labels[labels.WorkloadPreferredLocation] != "" {
		return "preferred"
	}

	return ""
}

func (g *location) parseLocations(locations string) {
	divisions := strings.Split(locations, "-")
	g.cities = strings.Split(divisions[0], "_")
	g.countries = strings.Split(divisions[1], "_")
	g.continents = strings.Split(divisions[2], "_")
}

func (g *location) findCountry(countryID string) (gountries.Country, error) {
	if country, err := g.query.FindCountryByName(countryID); err == nil {
		return country, nil
	}

	if country, err := g.query.FindCountryByAlpha(countryID); err == nil {
		return country, nil
	}

	return gountries.Country{}, errors.New("given country identifier does not match any country")
}

func getKeys(kvs map[string]bool) []string {
	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}
	return keys
}
