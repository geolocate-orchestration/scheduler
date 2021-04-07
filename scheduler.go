package scheduler

import (
	"errors"
	"scheduler/algorithms"
	"scheduler/algorithms/location"
	"scheduler/algorithms/naivelocation"
	"scheduler/algorithms/random"
	"scheduler/nodes"
)

// NewScheduler create a new instance of the IScheduler interface
func NewScheduler(algorithm string) (IScheduler, error) {
	s := &Scheduler{}

	if !algorithmExists(algorithm) {
		return nil, errors.New("selected algorithm does not exist")
	}

	s.inodes = nodes.New()
	s.initAlgorithm(algorithm)

	return s, nil
}

// ScheduleWorkload select a node based on used algorithm for the given workload
func (s *Scheduler) ScheduleWorkload(workload *algorithms.Workload) (*nodes.Node, error) {
	return s.algorithm.GetNode(workload)
}

// AddNode adds information about a new cluster node to the algorithm
func (s *Scheduler) AddNode(node *nodes.Node) {
	s.inodes.AddNode(node)
}

// UpdateNode updates information about a cluster node in the algorithm
func (s *Scheduler) UpdateNode(oldNode *nodes.Node, newNode *nodes.Node) {
	s.inodes.UpdateNode(oldNode, newNode)
}

// DeleteNode deletes information about a cluster node from the algorithm
func (s *Scheduler) DeleteNode(node *nodes.Node) {
	s.inodes.DeleteNode(node)
}


// Unexported

func algorithmExists(algorithmName string) bool {
	for _, algorithm := range AvailableAlgorithms {
		if algorithm == algorithmName {
			return true
		}
	}

	return false
}

func (s *Scheduler) initAlgorithm(algorithmName string) algorithms.Algorithm {
	var algorithm algorithms.Algorithm

	switch algorithmName {
	case "random":
		algorithm = random.New(s.inodes)
	case "naivelocation":
		algorithm = naivelocation.New(s.inodes)
	case "location":
		algorithm = location.New(s.inodes)
	default:
		algorithm = random.New(s.inodes)
	}

	return algorithm
}
