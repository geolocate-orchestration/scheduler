package scheduler

import (
	"github.com/geolocate-orchestration/scheduler/algorithms"
	"github.com/geolocate-orchestration/scheduler/nodes"
)

// AvailableAlgorithms list all package algorithms that can be used
var AvailableAlgorithms = []string{"location", "naivelocation", "random"}

// IScheduler exports all scheduler public methods
type IScheduler interface {
	// ScheduleWorkload returns a node selected from the chosen algorithm to bind the workload
	ScheduleWorkload(workload *algorithms.Workload) (*nodes.Node, error)

	// AddNode inserts new possible Node in the algorithm
	AddNode(node *nodes.Node)

	// UpdateNode replaces Node information in the algorithm
	UpdateNode(oldNode *nodes.Node, newNode *nodes.Node)

	// DeleteNode removes Node from the algorithm
	DeleteNode(node *nodes.Node)
}

// Scheduler has algorithm information
type Scheduler struct {
	inodes    nodes.INodes
	algorithm algorithms.Algorithm
}
