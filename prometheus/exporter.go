package prometheus

type Exporter struct {
	Endpoint string
	Path     string
}

func NewExporter(endpoint, path string) Exporter {
	return Exporter{
		Endpoint: endpoint,
		Path:     path,
	}
}
