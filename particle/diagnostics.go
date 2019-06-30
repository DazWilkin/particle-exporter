package particle

import (
	"fmt"
	"time"

	"github.com/DazWilkin/particle-exporter/prometheus"
)

const (
	urlDiagnostics = "https://api.particle.io/v1/diagnostics"
)

type DiagnosticsResponse struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}
type Diagnostic struct {
	UpdatedAt time.Time `json:"update_at"`
	Payload   Payload   `json:"payload"`
	ID        string    `json:"deviceID"`
}
type Payload struct {
	// Service Service `json:"service"`
	Device PDevice `json:"device"`
}
type Service struct {
}
type PDevice struct {
	System System `json:"system"`
}
type System struct {
	Memory Memory `json:"memory"`
	Uptime int64  `json:"uptime"`
}
type Memory struct {
	Used  int64 `json:"used"`
	Total int64 `json:"total"`
}

func (d Diagnostic) Expose() (result string) {
	result += fmt.Sprintf("# HELP device system memory used.\n# TYPE device_system_memory_used gauge\ndevice_system_memory_used{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Memory.Used)
	result += fmt.Sprintf("# HELP device system memory total.\n# TYPE device_system_memory_total gauge\ndevice_system_memory_total{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Memory.Total)
	result += fmt.Sprintf("# HELP device uptime.\n#TYPE device_uptime gauge\ndevice_uptime{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Uptime)
	return result
}
func (d Diagnostic) Export() prometheus.Gauge {
	ll := map[string]string{
		"device": d.ID,
	}
	return prometheus.NewGauge("device_system_memory_used", "Particle device system memory used", ll)
}
