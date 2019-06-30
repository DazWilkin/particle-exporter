package x

// Metric defines an interface for things that render in Prometheus exposition format
type Metric interface {
	Expose() string
}

// Producer defines a metrics source that can enqueue Metrics
type Producer interface {
	GetMetrics(chan Metric)
}

// Consumer defines a metrics sink that can dequeue Metrics
type Consumer interface {
}

// type Metric2 interface {
// 	Export() []prometheus.Exposer
// }
