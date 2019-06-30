package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/DazWilkin/particle-exporter/particle"
	"github.com/DazWilkin/particle-exporter/prometheus"
	"github.com/DazWilkin/particle-exporter/x"
)

var (
	endpoint = flag.String("endpoint", ":9999", "Endpoint for the Particle Exporter")
	path     = flag.String("path", "metrics", "Path on which Exporter should serve metrics")
	token    = flag.String("token", "", "Particle Access Token")
)
var (
	client particle.Client
)
var (
	metrics chan x.Metric
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[handler] Entering")
	var wgPut, wgGet sync.WaitGroup

	fmt.Fprintf(w, "# Particle Exporter")

	// Accumulate Metrics
	// Accumulate Devices|Diagnostics
	wgPut.Add(1)
	go func() {
		log.Println("[handler:go] Entered: Devices|Diagnostics")
		defer func() {
			log.Println("[handler:go] Exited: Devices|Diagnostics")
			wgPut.Done()
		}()
		log.Println("[handler] Getting Devices")
		devices, err := client.GetDevices()
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
			response, err := client.GetDiagnostics(device.ID)
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
	wgPut.Add(1)
	go func() {
		log.Println("[handler:go] Entered: Integrations")
		defer func() {
			log.Println("[handler:go] Exited: Integrations")
			wgPut.Done()
		}()
		log.Println("[handler] Getting Integrations")
		integrations, err := client.GetIntegrations()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[handler] Enumerating Integrations")
		for _, integration := range integrations {
			// Integration
			log.Printf("[handler] Getting Integration: %s", integration.ID)
			detailed, err := client.GetIntegration(integration.ID)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("[handler] Enqueuing Integration")
			metrics <- detailed
		}
	}()

	// Consume Metrics
	wgGet.Add(1)
	go func() {
		log.Println("[handler:go] Entered: Enumerate")
		defer func() {
			log.Println("[handler:go] Exited: Enumerate")
			wgGet.Done()
		}()
		log.Println("[handler] Enumerating Metrics")
		for metric := range metrics {
			log.Println("[handler] Metric")
			fmt.Fprintf(w, metric.Expose())
		}
	}()

	// Wait for Metric accumulation to complete
	// Then close the channel
	wgPut.Wait()
	close(metrics)

	// Wait for Metric consumption to complete
	// Then page is rendered
	wgGet.Wait()
}
func main() {
	log.Println("[main] Entered")
	flag.Parse()

	if *token == "" {
		log.Fatal("Require Particle Access Token to connect to Particle service.")
	}

	// Create Channel used to queue Metrics
	metrics = make(chan x.Metric)

	// Refactoring
	client = particle.NewClient(*token)
	_ = prometheus.NewExporter(*endpoint, *path)

	// Handle request for metrics
	log.Println("[main] Registering handler")
	http.HandleFunc(fmt.Sprintf("/%s", *path), metricsHandler)

	// Serve
	log.Printf("[main] Starting Server: %s", *endpoint)
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}
