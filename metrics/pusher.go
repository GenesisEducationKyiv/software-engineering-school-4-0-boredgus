package metrics

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultMetricsPushInterval time.Duration = 15 * time.Second

type (
	PushParams struct {
		URLToFetchMetrics string
		URLToPushMetrics  string
		PushInterval      time.Duration
	}

	Logger interface {
		Error(...any)
	}

	metricsPusher struct {
		logger Logger
	}
)

func NewMetricsPusher(logger Logger) *metricsPusher {
	return &metricsPusher{logger: logger}
}

func (p *metricsPusher) Push(ctx context.Context, params PushParams) {
	ticker := time.NewTicker(params.PushInterval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			p.push(params.URLToFetchMetrics, params.URLToPushMetrics)

			return

		case <-ticker.C:
			p.push(params.URLToFetchMetrics, params.URLToPushMetrics)
		}
	}
}

func (p *metricsPusher) push(from, to string) {
	data, err := p.getMetrics(from)
	if err != nil {
		p.logger.Error(err)

		return
	}

	if err = p.pushMetrics(to, data); err != nil {
		p.logger.Error(err)
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
	request, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return response.Body, nil
}
