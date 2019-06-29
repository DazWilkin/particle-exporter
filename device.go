package main

import (
	"fmt"
	"log"
)

type Device struct {
	ID                 string     `json:"id"`
	Name               string     `json:"string"`
	LastHeard          string     `json:"last_heard"`
	Connected          bool       `json:"connected"`
	Status             string     `json:"status"`
	SerialNumber       string     `json:"serial_number"`
	CurrentBuildTarget string     `"json:"current_build_target"`
	Variables          []Variable `json:"variables,omitempty"`
	Functions          []Function `json:"functions,omitempty"`
}
type Variable map[string]string
type Function []string

func (d Device) expose() string {
	log.Println("[device.expose] Entered")
	return fmt.Sprintf("# HELP particle device information.\n# TYPE particle_connected counter\nparticle_connected{core_id=\"%s\"} 1\n", d.ID)
}
