package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func SendRequest(ctx context.Context, url string, http_client http.Client) ([]byte, error) {

	var limiter = rate.NewLimiter(rate.Every(2*time.Second), 5)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {

		logger.Errorf("error creating request: %s", err)
		return nil, err
	}

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return nil, err
	}
	resp, err := http_client.Do(req)
	if err != nil {
		logger.Errorf("error response: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("unexpected status %d", resp.StatusCode)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error reading response body %s", err)
		return nil, err
	}
	return body, nil
}
