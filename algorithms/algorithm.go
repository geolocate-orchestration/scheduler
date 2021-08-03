package algorithms

import (
	"github.com/geolocate-orchestration/scheduler/nodes"
)

// Algorithm interface that exposes GetNode method
type Algorithm interface {
	GetName() string
	GetNode(pod *Workload) (*nodes.Node, error)
}

// Workload represents a cluster application to be scheduled
type Workload struct {
	// Name represents Workload unique identifying name
	Name string

	// Labels represents all labels Node must have to be considered for this workload schedule
	Labels map[string]string

	// CPU represents Workloads' necessary CPU resources Nodes must at least have available
	CPU int64

	// Memory represents Workloads' necessary Memory resources Nodes must at least have available
	Memory int64
}
