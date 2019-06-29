package main

import (
	"encoding/json"
)

const (
	urlDevices = "https://api.particle.io/v1/devices"
)

type Devices []Device

func newDevices(token string) (devices Devices, err error) {
	body, err := get(urlDevices, token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
func (d Devices) expose() (result []string) {
	result = make([]string, len(d))
	for i, device := range d {
		result[i] = device.expose()
	}
	return result
}
