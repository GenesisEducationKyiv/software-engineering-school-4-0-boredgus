package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsServer struct {
	logger Logger
	server *http.Server
}

func NewMetricsServer(logger Logger, httpAddr, route string, registry *prometheus.Registry) *metricsServer {
	mux := http.NewServeMux()
	mux.Handle(route, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	server := (&http.Server{
		Addr:    httpAddr,
		Handler: mux,
	})

	return &metricsServer{
		logger: logger,
		server: server,
	}
}

func (s *metricsServer) Run(ctx context.Context) {
	ch := make(chan struct{})
	ch <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			if err := s.server.Shutdown(context.Background()); err != nil {
				s.logger.Error("failed to shutdown metrics server: ", err)
			}

			return

		case <-ch:
			if err := s.server.ListenAndServe(); err != nil {
				s.logger.Error("failed to serve metrics server: ", err)
			}

			return
		}
	}
}
