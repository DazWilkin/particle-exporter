package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	urlIntegrations = "https://api.particle.io/v1/integrations"
)

type Integrations []Integration
type IntegrationResponse struct {
	Integration Integration `json:"integration"`
}
type Integration struct {
	ID                 string    `json:"id"`
	DeviceID           string    `json:"deviceID,omitempty"`
	Event              string    `json:"event"`
	CreatedAt          time.Time `json:"created_at"`
	Form               Form      `json:"form,omitempty"`
	Logs               []Log     `json:"logs,omitempty"`
	IntegrationType    string    `json:"integration_type"`
	URL                string    `json:"url"`
	RequestType        string    `json:"requestType"`
	IsJSON             bool      `json:"json"`
	NoDefaults         bool      `json:"noDefaults"`
	RejectUnauthorized bool      `json:"rejectUnauthorized"`
	Errors             []IError  `json:"errors,omitempty"`
	Counters           []Counter `json:"counters,omitempty"`
}
type Form struct {
	Field1 string `json:"field1"`
	APIKey string `json:"api_key"`
}
type Log struct {
	Event    Event     `json:"event"`
	Type     string    `json:"type"`
	Request  string    `json:"request"`
	Response string    `json:"response"`
	Time     time.Time `json:"time"`
}
type Event struct {
	Name        string    `json:"name"`
	Data        string    `json:"data"`
	TTL         int32     `json:"ttl"`
	PublishedAt time.Time `json:"published_at"`
	CoreID      string    `json:"coreid"`
}
type IError struct {
	Event    Event     `json:"event"`
	Type     string    `json:"type"`
	Request  string    `json:"request"`
	Response string    `json:"response"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}
type Counter struct {
	Date    string `json:"date"`
	Success string `json:"success"`
	Error   string `json:"error"`
}

func newIntegrations(token string) (Integrations, error) {
	body, err := get(urlIntegrations, token)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	ii := Integrations{}
	json.Unmarshal(body, &ii)
	log.Println(len(ii))
	return ii, nil
}
func newIntegration(token, id string) (Integration, error) {
	body, err := get(fmt.Sprintf("%s/%s", urlIntegrations, id), token)
	if err != nil {
		return Integration{}, err
	}
	log.Println(string(body))

	ir := IntegrationResponse{}
	ir.Integration.Logs = []Log{}
	ir.Integration.Errors = []IError{}
	ir.Integration.Counters = []Counter{}

	json.Unmarshal(body, &ir)
	return ir.Integration, nil
}

func (i Integration) expose() (result string) {
	log.Println("[integration.expose] Entered")
	result += fmt.Sprintf("# HELP particle integration logs information.\n# TYPE particle_integration_logs_count counter\nparticle_integration_logs_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Logs))
	result += fmt.Sprintf("# HELP particle integration errors information.\n# TYPE particle_integration_errors_count counter\nparticle_integration_errors_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Errors))
	result += fmt.Sprintf("# HELP particle integration counters nformation.\n# TYPE particle_integration_counters_count counter\nparticle_integration_counters_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Counters))
	return result
}
