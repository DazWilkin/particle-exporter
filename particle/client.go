package particle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DazWilkin/particle-exporter/x"
)

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) Client {
	return Client{
		client: &http.Client{},
		token:  token,
	}
}
func (c *Client) get(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (c *Client) GetDevices() (devices Devices, err error) {
	body, err := c.get(urlDevices)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
func (c *Client) GetDiagnostics(id string) (DiagnosticsResponse, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", urlDiagnostics, id))
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	dr := DiagnosticsResponse{}
	dr.Diagnostics = []Diagnostic{}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		return DiagnosticsResponse{}, err
	}
	return dr, nil
}
func (c *Client) GetIntegrations() (Integrations, error) {
	body, err := c.get(urlIntegrations)
	if err != nil {
		return nil, err
	}
	ii := Integrations{}
	json.Unmarshal(body, &ii)
	log.Println(len(ii))
	return ii, nil
}
func (c *Client) GetIntegration(id string) (Integration, error) {
	body, err := c.get(fmt.Sprintf("%s/%s", urlIntegrations, id))
	if err != nil {
		return Integration{}, err
	}

	ir := IntegrationResponse{}
	ir.Integration.Logs = []Log{}
	ir.Integration.Errors = []IError{}
	ir.Integration.Counters = []Counter{}

	json.Unmarshal(body, &ir)
	return ir.Integration, nil
}

// (TODO:dazwilkin) Refactor
type gauge struct {
	name  string
	group string
	help  string
	value time.Duration
}

func (g gauge) Expose() (result string) {
	result += fmt.Sprintf("# HELP %s.\n", g.help)
	result += fmt.Sprintf("# TYPE %s gauge\n", g.name)
	result += fmt.Sprintf("%s{type=\"%s\"} %d\n", g.name, g.group, g.value/time.Millisecond)
	return result
}
func (c Client) GetMetrics(metrics chan x.Metric) {
	var wg sync.WaitGroup

	// Accumulate Devices|Diagnostics
	wg.Add(1)
	go func() {
		startDevices := time.Now()
		log.Println("[handler:go] Entered: Devices|Diagnostics")
		defer func() {
			log.Println("[handler:go] Exited: Devices|Diagnostics")
			// Yes, another metric: measuring the elasped time
			metrics <- gauge{
				name:  "exporter_produce_time",
				group: "devices",
				help:  "Exporter Milliseconds to produce Metrics",
				value: time.Since(startDevices),
			}
			wg.Done()
		}()
		log.Println("[handler] Getting Devices")
		devices, err := c.GetDevices()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[handler] Enumerating Devices")
		for _, device := range devices {
			// Device
			log.Printf("[handler] Enqueuing Device: %s", device.ID)
			metrics <- device

			// Diagnostics
			log.Println("[handler] Getting Device Diagnostics")
			response, err := c.GetDiagnostics(device.ID)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("[handler] Enumerating Diagnostics")
			for _, diagnostic := range response.Diagnostics {
				log.Println("[handler] Enqueuing Diagnostics")
				metrics <- diagnostic
			}
		}
	}()

	// Accumulate Integrations
	wg.Add(1)
	go func() {
		startIntegrations := time.Now()
		log.Println("[handler:go] Entered: Integrations")
		defer func() {
			log.Println("[handler:go] Exited: Integrations")
			// Yes, another metric: measuring the elasped time
			metrics <- gauge{
				name:  "exporter_produce_time",
				group: "integrations",
				help:  "Exporter Milliseconds to produce Metrics",
				value: time.Since(startIntegrations),
			}
			wg.Done()
		}()
		log.Println("[handler] Getting Integrations")
		integrations, err := c.GetIntegrations()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[handler] Enumerating Integrations")
		for _, integration := range integrations {
			// Integration
			log.Printf("[handler] Getting Integration: %s", integration.ID)
			detailed, err := c.GetIntegration(integration.ID)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("[handler] Enqueuing Integration")
			log.Println(detailed.ID)
			metrics <- detailed
		}
	}()

	wg.Wait()
	// Done writing to channel, so close it
	close(metrics)
}
