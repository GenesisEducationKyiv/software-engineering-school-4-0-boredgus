package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExposeMetrics(httpAddr, route string, registry *prometheus.Registry) {
	mux := http.NewServeMux()
	mux.Handle(route, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	httpSrv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := httpSrv.ListenAndServe(); err != nil {
		fmt.Println("failed to server")
		return
	}

	fmt.Println("after listen and server")
}
