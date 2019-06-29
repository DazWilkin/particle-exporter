package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	endpoint = flag.String("endpoint", ":9999", "Endpoint for the Particle Exporter")
	path     = flag.String("path", "metrics", "Path on which Exporter should serve metrics")
	token    = flag.String("token", "", "Particle Access Token")
)
var (
	metrics []Metric
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "# Particle Exporter")

	// Enumerate metrics
	metrics = append(metrics, newFreddie())

	// Enumerate Devices|Diagnostics
	devices, err := newDevices(*token)
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		// Device
		metrics = append(metrics, device)

		// Diagnostics
		response, err := newDiagnostics(*token, device.ID)
		if err != nil {
			log.Fatal(err)
		}
		for _, diag := range response.Diagnostics {
			metrics = append(metrics, diag)
		}
	}

	// Enumerate Integrations
	integrations, err := newIntegrations(*token)
	if err != nil {
		log.Fatal(err)
	}
	for _, integration := range integrations {
		// Integration
		detailed, err := newIntegration(*token, integration.ID)
		if err != nil {
			log.Fatal(err)
		}
		metrics = append(metrics, detailed)
	}

	for _, metric := range metrics {
		fmt.Fprintf(w, metric.expose())
	}
}
func main() {
	flag.Parse()

	if *token == "" {
		log.Fatal("Require Particle Access Token to connect to Particle service.")
	}

	// Handle request for metrics
	http.HandleFunc(fmt.Sprintf("/%s", *path), metricsHandler)
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}
