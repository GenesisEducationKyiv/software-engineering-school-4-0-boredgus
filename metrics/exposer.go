package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExposeMetrics(httpAddr, route string, registry *prometheus.Registry) error {
	mux := http.NewServeMux()
	mux.Handle(route, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	return (&http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}).ListenAndServe()
}
