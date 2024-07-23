package metrics

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/gateway/internal/config/logger"
)

const (
	MetricsPushInterval time.Duration = 15 * time.Second
	requestTimeout      time.Duration = 2 * time.Second
)

type metricsPusher struct {
	logger logger.Logger
}

func NewMetricsPusher(logger logger.Logger) *metricsPusher {
	return &metricsPusher{logger: logger}
}

func (p *metricsPusher) Push(urlToFetchMetrics, urlToPushMetrics string, pushInterval time.Duration) {
	ticker := time.NewTicker(pushInterval)

	for range ticker.C {
		data, err := p.getMetrics(urlToFetchMetrics)
		if err != nil {
			p.logger.Error(err)

			continue
		}

		if err = p.pushMetrics(urlToPushMetrics, data); err != nil {
			p.logger.Error(err)

			continue
		}
	}
}

func (p *metricsPusher) pushMetrics(url string, data io.Reader) error {
	body, err := p.makeRequest(http.MethodPost, url, data)
	if err != nil {
		return fmt.Errorf("failed to push metrics: %w", err)
	}
	defer body.Close()

	return nil
}

func (p *metricsPusher) getMetrics(url string) (io.Reader, error) {
	body, err := p.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer body.Close()

	var buffer bytes.Buffer
	if _, err = buffer.ReadFrom(body); err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return &buffer, nil
}

func (p *metricsPusher) makeRequest(method, url string, body io.Reader) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return response.Body, nil
}
