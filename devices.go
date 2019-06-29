package main

import (
	"encoding/json"
	"log"
)

const (
	url = "https://api.particle.io/v1/devices"
)

type Devices []Device

func newDevices(token string) (devices Devices, err error) {
	log.Println("[devices:new] Entered")
	body, err := get(url, token)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &devices)
	return devices, nil
}
func (d Devices) expose() (result []string) {
	log.Println("[devices.expose] Entered")
	result = make([]string, len(d))
	for i, device := range d {
		result[i] = device.expose()
	}
	return result
}
