package labels

// Node indicates needed Node label for it to be considered in the algorithms
const Node = "node.geolocate.io"

// NodeCity indicates Node city location
const NodeCity = "node.geolocate.io/city"

// NodeCountry indicates Node country location
const NodeCountry = "node.geolocate.io/country"

// NodeContinent indicates Node continent location
const NodeContinent = "node.geolocate.io/continent"

// WorkloadRequiredLocation indicates Workloads required Node location to be scheduled there
const WorkloadRequiredLocation = "workload.geolocate.io/requiredLocation"

// WorkloadPreferredLocation indicates Workloads preferred Node location to be prioritized in scheduling
const WorkloadPreferredLocation = "workload.geolocate.io/preferredLocation"
