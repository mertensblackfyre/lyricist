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

	simple, complex := BuildQueryAPI(track)
	complex_body, err := SendRequest(ctx, complex, http.Client{
		Timeout: 15 * time.Second,
	})

	if err != nil {
		return TrackDeezer{}, err
	}

	var result DeezerSearchResponse
	if err := json.Unmarshal(complex_body, &result); err != nil {

		logger.Errorf("unmarshaling error of complex query: %s", err)
		return TrackDeezer{}, err
	}
	if len(result.Data) == 0 {
		logger.Warn("no data using complex query,trying simple query")

		simple_body, err := SendRequest(ctx, simple, http.Client{
			Timeout: 15 * time.Second,
		})

		if err != nil {
			return TrackDeezer{}, err
		}

		var result DeezerSearchResponse
		if err := json.Unmarshal(simple_body, &result); err != nil {
			logger.Errorf("unmarshaling error of simple query: %s", err)
			return TrackDeezer{}, err
		}
		if len(result.Data) == 0 {
			logger.Errorf("no data found using both queries")
			return TrackDeezer{}, err
		}

		t := result.Data[0]
		return t, nil
	}

	t := result.Data[0]
	return t, nil
}

func BuildQueryAPI(track *TrackScrapeInfo) (string, string) {

	q := fmt.Sprintf(`"%s"`, CleanTrackTitle(track.Title))
	params := url.Values{}
	params.Set("q", q)
	simple_url := base + "search?" + params.Encode()

	q = fmt.Sprintf(`artist:"%s" track:"%s"`, track.Artists[0], CleanTrackTitle(track.Title))
	params = url.Values{}
	params.Set("q", q)

	complex_url := base + "search?" + params.Encode()

	return simple_url, complex_url

}

func CleanTrackTitle(title string) string {
	reParen := regexp.MustCompile(`(?i)\s*\(.*?(feat|ft|featuring).*?\)`)
	title = reParen.ReplaceAllString(title, "")

	re_bracket := regexp.MustCompile(`(?i)\s*\[.*?(official (audio|video|music video)|lyrics|hd|clean|explicit).*?\]`)
	title = re_bracket.ReplaceAllString(title, "")

	re_suffix := regexp.MustCompile(`(?i)\s*[-–]\s*(single|remaster(ed)?|deluxe edition|album version).*$`)
	title = re_suffix.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}
