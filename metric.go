package main

// Metric defines an interface for things that render in Prometheus exposition format
type Metric interface {
	expose() string
}
