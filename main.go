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
	fmt.Fprintf(w, "# Hello Freddie")
	// Enumerate metrics
	for _, metric := range metrics {
		fmt.Fprintf(w, metric.expose())
	}
}
func main() {
	flag.Parse()

	if *token == "" {
		log.Fatal("Require Particle Access Token to connect to Particle service.")
	}

	// Enumerate metrics
	metrics = append(metrics, newFreddie())

	// Devices
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

	// Handle request for metrics
	http.HandleFunc(fmt.Sprintf("/%s", *path), metricsHandler)
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}
