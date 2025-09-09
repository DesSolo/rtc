package client

import (
	"context"
	"fmt"
	"net/http"
)

// UpsertConfigRequest ...
type UpsertConfigRequest struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"value_type"`
	Usage     string `json:"usage"`
	Group     string `json:"group"`
	Writable  bool   `json:"writable"`
}

// UpsertConfigs ...
func (c *Client) UpsertConfigs(ctx context.Context, projectName, envName, releaseName string, req []*UpsertConfigRequest) error {
	uri := fmt.Sprintf("/projects/%s/envs/%s/releases/%s/configs", projectName, envName, releaseName)

	body, err := encodePayload(req)
	if err != nil {
		return fmt.Errorf("marshalling upsert configs: %w", err)
	}

	httpReq, err := c.newRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return fmt.Errorf("newRequest: %w", err)
	}

	if _, err := c.do(httpReq, http.StatusCreated); err != nil {
		return fmt.Errorf("do: %w", err)
	}

	return nil
}
