package algorithms

import (
	"scheduler/nodes"
)

// Algorithm interface that exposes GetNode method
type Algorithm interface {
	GetName() string
	GetNode(pod *Workload) (*nodes.Node, error)
}

// Workload represents a cluster application to be scheduled
type Workload struct {
	Name   string
	Labels map[string]string
	CPU    int64
	Memory int64
}
