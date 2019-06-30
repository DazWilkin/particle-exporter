package particle

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDiagnostic(t *testing.T) {
	body, err := get(urlDiagnostics, tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	dr := DiagnosticsResponse{}
	dr.Diagnostics = []Diagnostic{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}

func TestDiagnostics(t *testing.T) {
	body, err := get(fmt.Sprintf("%s/%s", urlDiagnostics, device), tk)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
	dr := DiagnosticsResponse{}
	dr.Diagnostics = []Diagnostic{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		t.Errorf("Problem: %s", err)
	}
}
