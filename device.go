package main

import (
	"fmt"
	"time"
)

type Device struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	LastApp            string     `json:"last_app"`
	LastIPAddress      string     `json:"last_ip_address"`
	LastHeard          time.Time  `json:"last_heard"`
	ProductID          int32      `json:"product_id"`
	Connected          bool       `json:"connected"`
	PlatformID         int32      `json:"platform_id"`
	Cellular           bool       `json:"cellular"`
	Notes              []string   `json:"notes"`
	Status             string     `json:"status"`
	SerialNumber       string     `json:"serial_number"`
	MobileSecret       string     `json:"mobile_secret"`
	CurrentBuildTarget string     `json:"current_build_target"`
	SystemFirmware     string     `json:"system_firmware_version"`
	PinnedBuildTarget  string     `json:"pinned_build_target"`
	Variables          []Variable `json:"variables,omitempty"`
	Functions          []Function `json:"functions,omitempty"`
}
type Variable map[string]string
type Function []string

func (d Device) expose() string {
	return fmt.Sprintf("# HELP particle device information.\n# TYPE particle_connected counter\nparticle_connected{core_id=\"%s\"} 1\n", d.ID)
}
