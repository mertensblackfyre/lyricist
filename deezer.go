package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var base = "https://api.deezer.com/"

func GetTrack(ctx context.Context, track *TrackScrapeInfo) (TrackDeezer, error) {

	q := fmt.Sprintf(`artist:"%s" track:"%s"`, track.Artists[0], CleanTrackTitle(track.Title))
	params := url.Values{}
	params.Set("q", q)

	url := base + "search?" + params.Encode()

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackDeezer{}, fmt.Errorf("creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return TrackDeezer{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TrackDeezer{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return TrackDeezer{}, err
	}

	type DeezerSearchResponse struct{ Data []TrackDeezer }
	var searchResp DeezerSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return TrackDeezer{}, fmt.Errorf("unmarshaling error: %w", err)
	}
	if len(searchResp.Data) == 0 {
		return TrackDeezer{}, fmt.Errorf("track name missing in JSON-LD")
	}

	t := searchResp.Data[0]
	return t, nil
}

func CleanTrackTitle(title string) string {
	reParen := regexp.MustCompile(`(?i)\s*\(.*?(feat|ft|featuring).*?\)`)
	title = reParen.ReplaceAllString(title, "")

	// Remove common bracketed tags
	reBracket := regexp.MustCompile(`(?i)\s*\[.*?(official (audio|video|music video)|lyrics|hd|clean|explicit).*?\]`)
	title = reBracket.ReplaceAllString(title, "")

	// Remove trailing " - Single", " - Remastered", etc. (optional)
	reSuffix := regexp.MustCompile(`(?i)\s*[-–]\s*(single|remaster(ed)?|deluxe edition|album version).*$`)
	title = reSuffix.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}
