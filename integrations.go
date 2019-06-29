package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	integs = "https://api.particle.io/v1/integrations"
)

type Integrations []Integration
type Integration struct {
	ID                 string    `json:"id"`
	DeviceID           string    `json:"deviceID,omitempty"`
	Event              string    `json:"event"`
	CreatedAt          time.Time `json:"created_at"`
	Logs               []Log     `json:"logs,omitempty"`
	IntegrationType    string    `json:"integration_type"`
	URL                string    `json:"url"`
	RequestType        string    `json:"requestType"`
	IsJSON             bool      `json:"json"`
	NoDefaults         bool      `json:"noDefaults"`
	RejectUnauthorized bool      `json:"rejectUnauthorized"`
	Errors             []IError  `json:"errors,omitempty"`
	Counters           []Counter `json:"counters"`
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
type IError struct{}
type Counter struct {
	Date    string `json:"date"`
	Success string `json:"success"`
}

func newIntegrations(token string) (integrations Integrations, err error) {
	log.Println("[integrations:new] Entered")
	body, err := get(integs, token)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &integrations)
	return integrations, nil
}
func newIntegration(token, id string) (integration Integration, err error) {
	log.Println("[integration:new] Entered")
	body, err := get(fmt.Sprintf("%s/%s", integs, id), token)
	if err != nil {
		return Integration{}, err
	}
	json.Unmarshal(body, &integration)
	return integration, nil
}

func (i Integration) expose() (result string) {
	log.Println("[integration.expose] Entered")
	result += fmt.Sprintf("# HELP particle integration logs information.\n# TYPE particle_integration_logs_count counter\nparticle_integration_logs_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Logs))
	result += fmt.Sprintf("# HELP particle integration errors information.\n# TYPE particle_integration_errors_count counter\nparticle_integration_errors_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Errors))
	result += fmt.Sprintf("# HELP particle integration counters nformation.\n# TYPE particle_integration_counters_count counter\nparticle_integration_counters_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Counters))
	return result
}
