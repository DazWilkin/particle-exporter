package prometheus

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DazWilkin/particle-exporter/x"
)

type Exporter struct {
	Endpoint string
	Path     string
	Producer x.Producer
}

// (TODO:dazwilkin) Refactor
type gauge struct {
	name  string
	help  string
	value time.Duration
}

func (g gauge) Expose() (result string) {
	result += fmt.Sprintf("# HELP %s.\n", g.help)
	result += fmt.Sprintf("# TYPE %s gauge\n", g.name)
	result += fmt.Sprintf("%s %d\n", g.name, g.value/time.Millisecond)
	return result
}
func NewExporter(endpoint, path string, producer x.Producer) Exporter {
	return Exporter{
		Endpoint: endpoint,
		Path:     path,
		Producer: producer,
	}
}
func (e *Exporter) metricHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[handler] Entering")

	// Add a headline
	fmt.Fprintf(w, "# Particle Exporter")

	// Create channel for Metrics
	metrics := make(chan x.Metric)

	// Consume Metrics
	// Important: must provide reader from channel *before* writing
	// Otherwise writer will block
	var wg sync.WaitGroup
	wg.Add(1)
	var elapsed time.Duration
	go func() {
		start := time.Now()
		log.Println("[handler:go] Entered: Enumerate")
		defer func() {
			log.Println("[handler:go] Exited: Enumerate")
			// Yes, another metric: measuring the elasped time
			// Too late to include this in the Consumer
			// So it's time until this defer is called
			// But executed before the handler (!) concludes
			elapsed = time.Since(start)
			wg.Done()
		}()
		log.Println("[handler] Enumerating Metrics")
		for metric := range metrics {
			log.Println("[handler] Metric")
			fmt.Fprintf(w, metric.Expose())
		}
	}()

	// Produce Metrics
	// Important: must provider writer to channel *after* reading
	e.Producer.GetMetrics(metrics)

	// Handler completes once Consumer completes
	wg.Wait()
	// Once Consumer completes, we can write out its elasped time
	fmt.Fprintf(w, gauge{
		name:  "exporter_consume_time",
		help:  "Exporter Milliseconds to consume Metrics",
		value: elapsed,
	}.Expose())
}
func (e *Exporter) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}
func (e *Exporter) Run() {
	log.Println("[exporter:run] Registering handler")
	http.HandleFunc(fmt.Sprintf("/%s", e.Path), e.metricHandler)
	http.HandleFunc("/healthz", e.healthHandler)
	log.Printf("[exporter:run] Starting Server: %s", e.Endpoint)
	log.Fatal(http.ListenAndServe(e.Endpoint, nil))
}
