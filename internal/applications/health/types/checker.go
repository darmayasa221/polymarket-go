// Package types defines types for the health application service.
package types

// ComponentStatus holds the health status of one system component.
type ComponentStatus struct {
	Name    string
	Status  string
	Message string
}

// HealthReport holds the overall health of the system.
type HealthReport struct {
	Status     string
	Components []ComponentStatus
}
