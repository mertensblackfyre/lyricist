package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var base = "https://api.deezer.com/"

func GetTrack(ctx context.Context, track *Track) (Track, error) {
	q := fmt.Sprintf(`artist:"%s" track:"%s"`, track.Artists[0], track.Title)
	params := url.Values{}
	params.Set("q", q)

	url := base + "search?" + params.Encode()

	client := &http.Client{
		Timeout: time.Second * 10, // Timeout each requests
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Track{}, fmt.Errorf("creating request: %w", err)
	}
	resp, err := client.Do(req)
	fmt.Println(resp)
	if err != nil {
		return Track{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Track{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Track{}, err
	}

	return *track, nil
}
