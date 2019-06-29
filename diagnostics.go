package main

import (
	"encoding/json"
	"fmt"
	"time"
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

func newDiagnostics(token, device string) (DiagnosticsResponse, error) {
	body, err := get(fmt.Sprintf("%s/%s", urlDiagnostics, device), token)
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	dr := DiagnosticsResponse{}
	dr.Diagnostics = []Diagnostic{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	return dr, nil
}
func (d Diagnostic) expose() (result string) {
	result += fmt.Sprintf("# HELP device system memory used.\n# TYPE device_system_memory_used gauge\ndevice_system_memory_used{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Memory.Used)
	result += fmt.Sprintf("# HELP device system memory total.\n# TYPE device_system_memory_total gauge\ndevice_system_memory_total{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Memory.Total)
	result += fmt.Sprintf("# HELP device uptime.\n#TYPE device_uptime gauge\ndevice_uptime{core_id=\"%s\"} %d\n", d.ID, d.Payload.Device.System.Uptime)
	return result
}
