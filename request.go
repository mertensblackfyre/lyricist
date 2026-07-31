package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func SendRequest(ctx context.Context, url string, http_client http.Client) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := http_client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	return body, nil
}
