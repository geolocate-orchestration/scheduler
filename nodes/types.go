package nodes

import (
	"github.com/aida-dos/gountries"
)

// INodes exports all node controller public methods
type INodes interface {
	CountNodes() int
	GetAllNodes() []*Node
	GetNodes(filter *NodeFilter) []*Node

	AddNode(node *Node)
	UpdateNode(oldNode *Node, newNode *Node)
	DeleteNode(node *Node)
}

// New create a new Nodes struct
func New() INodes {
	nodes := Nodes{
		Query:          gountries.New(),
		ContinentsList: gountries.NewContinents(),

		Nodes:      make([]*Node, 0),
		Cities:     make(map[string][]*Node),
		Countries:  make(map[string][]*Node),
		Continents: make(map[string][]*Node),
	}

	return &nodes
}

// Nodes controls in-cache nodes
type Nodes struct {
	Query          *gountries.Query
	ContinentsList gountries.Continents

	Nodes      []*Node
	Cities     map[string][]*Node
	Countries  map[string][]*Node
	Continents map[string][]*Node
}

// Node represents a cluster Node
type Node struct {
	Name   string
	Labels map[string]string
	CPU    int64
	Memory int64
}

// NodeFilter states the params which nodes must match to be returned
type NodeFilter struct {
	Labels    map[string]string
	Resources Resources
	Locations Locations
}

// Resources states the available resources nodes must have to be returned
type Resources struct {
	CPU    int64
	Memory int64
}

// Locations states the location params nodes must match to be returned
type Locations struct {
	Cities     []string
	Countries  []string
	Continents []string
}

const (
	cityLabel      = "node.edge.aida.io/city"
	countryLabel   = "node.edge.aida.io/country"
	continentLabel = "node.edge.aida.io/continent"
)