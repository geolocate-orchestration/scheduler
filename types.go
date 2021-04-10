package scheduler

import (
	"github.com/mv-orchestration/scheduler/algorithms"
	"github.com/mv-orchestration/scheduler/nodes"
)

// AvailableAlgorithms list all package algorithms that can be used
var AvailableAlgorithms = []string{"location", "naivelocation", "random"}

// IScheduler exports all scheduler public methods
type IScheduler interface {
	ScheduleWorkload(workload *algorithms.Workload) (*nodes.Node, error)

	AddNode(node *nodes.Node)
	UpdateNode(oldNode *nodes.Node, newNode *nodes.Node)
	DeleteNode(node *nodes.Node)
}

// Scheduler has algorithm information
type Scheduler struct {
	inodes    nodes.INodes
	algorithm algorithms.Algorithm
}
