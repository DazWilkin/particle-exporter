package particle

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIntegrations(t *testing.T) {
	body, err := get(urlIntegrations, tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	integrations := Integrations{}
	err = json.Unmarshal(body, &integrations)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}
func TestIntegration(t *testing.T) {
	body, err := get(fmt.Sprintf("%s/%s", urlIntegrations, integration), tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	ir := IntegrationResponse{}
	err = json.Unmarshal(body, &ir)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}
