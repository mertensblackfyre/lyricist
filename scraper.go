package main

import (
	"context"
	"net/http"
	"regexp"
	"time"
)

func removeJapanese(s string) string {
	return regexp.MustCompile(`[^\x00-\x7F]+`).ReplaceAllString(s, "")
}
func ScrapeTrackInfo(ctx context.Context, url string) (TrackScrapeInfo, error) {

	body, err := SendRequest(ctx, url,
		http.Client{
			Timeout: 15 * time.Second,
		})

	if err != nil {
		return TrackScrapeInfo{}, err
	}

	final, err := ParseRecordingJSON(string(body))
	if err != nil {
		logger.Errorf("extracting JSON-LD: %s", err)
		return TrackScrapeInfo{}, err
	}

	art := ParseArtists(final[1])

	return TrackScrapeInfo{
		Title:   final[0],
		Artists: art,
	}, nil
}
