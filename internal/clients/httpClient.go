package clients

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var InvalidRequestErr = errors.New("invalid-request")

type HTTPClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

func NewHTTPClient() *httpClient {
	return &httpClient{
		client: &http.Client{},
	}
}

// Creates new http.Request.
func (c *httpClient) request(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", InvalidRequestErr, err)
	}

	return req, nil
}

// Sends http.Request and returns http.Response
func (c *httpClient) do(request *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", InvalidRequestErr, err)
	}

	return resp, nil
}

// Sends GET request with context.
func (c *httpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := c.request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}
