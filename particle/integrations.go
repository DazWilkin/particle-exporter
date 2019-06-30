package particle

import (
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
	ID        string `json:"id"`
	DeviceID  string `json:"deviceID,omitempty"`
	Event     string `json:"event"`
	CreatedAt string `json:"created_at"`
	// Form               Form      `json:"form,omitempty"`
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
	Event    Event  `json:"event"`
	Type     string `json:"type"`
	Request  string `json:"request"`
	Response string `json:"response"`
	Time     int64  `json:"time"`
}
type Event struct {
	Name        string    `json:"name"`
	Data        string    `json:"data"`
	TTL         int32     `json:"ttl"`
	PublishedAt time.Time `json:"published_at"`
	CoreID      string    `json:"coreid"`
}
type IError struct {
	Event    Event  `json:"event"`
	Type     string `json:"type"`
	Request  string `json:"request"`
	Response string `json:"response"`
	Message  string `json:"message"`
	Time     int64  `json:"time"`
}
type Counter struct {
	Date    string `json:"date"`
	Success string `json:"success"`
	Error   string `json:"error"`
}

func (i Integration) Expose() (result string) {
	log.Println("[integration.expose] Entered")
	result += fmt.Sprintf("# HELP particle integration logs information.\n# TYPE particle_integration_logs_count counter\nparticle_integration_logs_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Logs))
	result += fmt.Sprintf("# HELP particle integration errors information.\n# TYPE particle_integration_errors_count counter\nparticle_integration_errors_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Errors))
	result += fmt.Sprintf("# HELP particle integration counters nformation.\n# TYPE particle_integration_counters_count counter\nparticle_integration_counters_count{integration_id=\"%s\"} %d\n", i.ID, len(i.Counters))
	return result
}
