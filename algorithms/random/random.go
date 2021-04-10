package random

import (
	"errors"
	"github.com/mv-orchestration/scheduler/algorithms"
	"github.com/mv-orchestration/scheduler/nodes"
	"k8s.io/klog/v2"
)

type random struct {
	inodes nodes.INodes
}

// New creates new random struct
func New(inodes nodes.INodes) algorithms.Algorithm {
	return &random{
		inodes: inodes,
	}
}

func (r random) GetName() string {
	return "random"
}

func (r random) GetNode(*algorithms.Workload) (*nodes.Node, error) {
	klog.Infoln("getting cached nodes")
	return getRandomNode(r.inodes)
}

// GetRandomNode returns a random
func getRandomNode(inodes nodes.INodes) (*nodes.Node, error) {
	allNodes := inodes.GetAllNodes()

	if len(allNodes) == 0 {
		errMessage := "no nodes are available"
		return nil, errors.New(errMessage)
	}

	klog.Infof("will randomly get 1 node from the %d available\n", len(allNodes))
	return nodes.GetRandomFromList(allNodes)
}
