package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var base = "https://api.deezer.com/"

func GetTrackMetaData(ctx context.Context, track *TrackScrapeInfo) (TrackDeezer, error) {

	q := fmt.Sprintf(`artist:"%s" track:"%s"`, track.Artists[0], CleanTrackTitle(track.Title))
	params := url.Values{}
	params.Set("q", q)

	url := base + "search?" + params.Encode()

	body, err := SendRequest(ctx, url, http.Client{
		Timeout: 15 * time.Second,
	})

	if err != nil {
		return TrackDeezer{}, err
	}

	type Deezeronse struct{ Data []TrackDeezer }
	var search_resp Deezeronse
	if err := json.Unmarshal(body, &search_resp); err != nil {
		return TrackDeezer{}, fmt.Errorf("unmarshaling error: %w", err)
	}

	if len(search_resp.Data) == 0 {
		return TrackDeezer{}, fmt.Errorf("No data found")
	}

	t := search_resp.Data[0]
	return t, nil
}

func CleanTrackTitle(title string) string {
	reParen := regexp.MustCompile(`(?i)\s*\(.*?(feat|ft|featuring).*?\)`)
	title = reParen.ReplaceAllString(title, "")

	// Remove common bracketed tags
	re_bracket := regexp.MustCompile(`(?i)\s*\[.*?(official (audio|video|music video)|lyrics|hd|clean|explicit).*?\]`)
	title = re_bracket.ReplaceAllString(title, "")

	// Remove trailing " - Single", " - Remastered", etc. (optional)
	re_suffix := regexp.MustCompile(`(?i)\s*[-–]\s*(single|remaster(ed)?|deluxe edition|album version).*$`)
	title = re_suffix.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}
