package prometheus

import (
	"fmt"
)

type Exposer interface {
	String() string
}
type Base struct {
	Name   string
	Help   string
	Type   string
	Labels Labels
}
type Labels map[string]string

func (ll Labels) String() (result string) {
	for k, v := range ll {
		result += fmt.Sprintf("{%s:\"%s\"", k, v)
	}
	return result
}

type Counter struct {
	Base
	Value int64
}

func NewCounter(name, help string, ll Labels) Gauge {
	return Gauge{
		Base: Base{
			Name:   name,
			Help:   help,
			Type:   "counter",
			Labels: ll,
		},
	}
}
func (c Counter) Set(value int64) {
	c.Value = value
}
func (c Counter) String() (result string) {
	result += fmt.Sprintf("# HELP %s.\n", c.Help)
	result += fmt.Sprintf("# TYPE %s %s\n", c.Name, c.Type)
	result += fmt.Sprintf("%s{%s} %d\n", c.Name, c.Labels.String(), c.Value)
	return result
}

type Gauge struct {
	Base
	Value float64
}

func NewGauge(name, help string, ll Labels) Gauge {
	return Gauge{
		Base: Base{
			Name:   name,
			Help:   help,
			Type:   "gauge",
			Labels: ll,
		},
	}
}
func (g *Gauge) Set(value float64) {
	g.Value = value
}
func (g Gauge) String() (result string) {
	result += fmt.Sprintf("# HELP %s.\n", g.Help)
	result += fmt.Sprintf("# TYPE %s %s\n", g.Name, g.Type)
	result += fmt.Sprintf("%s{%s} %f\n", g.Name, g.Labels.String(), g.Value)
	return result
}
