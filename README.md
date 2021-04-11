# scheduler

[![Test](https://github.com/mv-orchestration/scheduler/actions/workflows/test.yml/badge.svg?branch=develop)](https://github.com/mv-orchestration/scheduler/actions/workflows/test.yml)

## Usage

### Scheduler

[types.go](types.go)
```go
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
```

### Nodes

[nodes/types.go](nodes/types.go)
```go
// Node represents a cluster Node
type Node struct {
	// Name represents Node unique identifying name
	Name string
	
	// Labels represents all of Node labels
	Labels map[string]string
	
	// CPU represents Node available CPU resources in MilliValue
	CPU int64

	// Memory represents Node available Memory resources in MilliValue
	Memory int64
}
```

Nodes can be configured with the following labels:

- **node.mv.io** - Node must have this label to be used in the algorithm
- **node.mv.io/city** - Indicates node city location
- **node.mv.io/country** - Indicates node country location
- **node.mv.io/continent** - Indicates node continent location

### Workload Labeling

[algorithms/algorithm.go](algorithms/algorithm.go)
```go
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
```

Workloads can be configured with the following labels:

- **workload.mv.io/requiredLocation** - List of Workload required locations
- **workload.mv.io/preferredLocation** - List of Workload preferred locations

Location format:

`(<CITY>?(_<CITY>)*)-(<COUNTRY>?(_<COUNTRY>)*))-(<CONTINENT>?(_<CONTINENT>)*))`

Examples:

`<CITY_A>-<COUNTRY_A>_<COUNTRY_B>-<CONTINENT_A>` -> `Braga-Portugal_Spain-Europe`

`<CITY_A>_<CITY_B>--<CONTINENT_A>` -> `Braga_Porto--Europe`

## Development

### Lint
```shell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
golangci-lint ./...
```

### Testing and Coverage
```shell
go test --coverprofile=coverage.out ./...
go tool cover -html=coverage.out 
```

### Format

```shell
go fmt ./...
```
