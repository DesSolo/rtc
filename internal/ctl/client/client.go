package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	url    string
	token  string
	client httpClient
}

// NewClient ...
func NewClient(url, token string) *Client {
	return &Client{
		url:    url,
		token:  token,
		client: http.DefaultClient,
	}
}

func (c *Client) newRequest(ctx context.Context, method, uri string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.url+uri, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))

	return req, nil
}

func (c *Client) do(req *http.Request, validStatusCode int) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	if resp.StatusCode != validStatusCode {
		data, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		return nil, fmt.Errorf("client.Do: invalid status code: %d message: %s", resp.StatusCode, data)
	}

	return resp, nil
}
