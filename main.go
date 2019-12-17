package main

import (
	"flag"
	"log"

	"github.com/DazWilkin/particle-exporter/particle"
	"github.com/DazWilkin/particle-exporter/prometheus"
)

var (
	endpoint = flag.String("endpoint", ":9375", "Endpoint for the Particle Exporter")
	path     = flag.String("path", "metrics", "Path on which Exporter should serve metrics")
	token    = flag.String("token", "", "Particle Access Token")
)

func main() {
	flag.Parse()

	if *token == "" {
		log.Fatal("Require Particle Access Token to connect to Particle service.")
	}

	// Create a client to Particle's Cloud
	client := particle.NewClient(*token)

	// Create a Prometheus Exporter client
	exporter := prometheus.NewExporter(*endpoint, *path, client)

	// Dequeue and render metrics
	exporter.Run()
}
