package x

import "github.com/DazWilkin/particle-exporter/prometheus"

// Metric defines an interface for things that render in Prometheus exposition format
type Metric interface {
	Expose() string
}
type Metric2 interface {
	Export() []prometheus.Exposer
}