package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ytdlp "github.com/lrstanley/go-ytdlp"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(3*time.Second), 2)
var dl = ytdlp.New()

func DownloadTrack(ctx context.Context, url string) (string, error) {
	pos := strings.Index(url, "v=")
	if pos == -1 {
		return "", fmt.Errorf("can't extract id from youtube url")
	}
	id := url[pos+2:]

	if err := limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit: %w", err)
	}
	_, err := dl.Run(ctx, url, "-f", "bestaudio", "-o", "%tempdir%/%(id)s.%(ext)s", "--no-playlist", "--update")
	if err != nil {
		return "", fmt.Errorf("error running ytdlp %w", err)

	}
	return id, nil
}
func Search(ctx context.Context, query string) (ytdlp.ExtractedInfo, error) {
	var builder strings.Builder

	builder.WriteString("ytsearch1:")
	builder.WriteString(query)

	result := builder.String()

	if err := limiter.Wait(ctx); err != nil {
		return ytdlp.ExtractedInfo{}, fmt.Errorf("rate limit: %w", err)
	}

	output, err := dl.Run(ctx, result, "--print-json", "--skip-download", "--no-playlist")
	if err != nil {
		return ytdlp.ExtractedInfo{}, fmt.Errorf("error running ytdlp %w", err)
	}
	var info ytdlp.ExtractedInfo
	if err := json.Unmarshal([]byte(output.Stdout), &info); err != nil {
		return ytdlp.ExtractedInfo{}, fmt.Errorf("error unmarshaling %w", err)
	}
	return info, nil
}

func BuildSearchQuery(track *TrackDeezer) string {

	var builder strings.Builder

	builder.WriteString(track.Artist.Name)
	builder.WriteString(", ")

	builder.WriteString(track.Title)
	builder.WriteString(" audio")
	url := builder.String()
	return url
}
