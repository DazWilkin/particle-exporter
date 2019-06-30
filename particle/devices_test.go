package particle

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDevices(t *testing.T) {
	body, err := get(urlDevices, tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	devices := Devices{}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}
func TestDevice(t *testing.T) {
	body, err := get(fmt.Sprintf("%s/%s", urlDevices, device), tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	device := Device{}
	err = json.Unmarshal(body, &device)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}
