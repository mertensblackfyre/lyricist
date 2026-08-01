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
		logger.Errorf("can't extract id from youtube url")
		return "", nil
	}
	id := url[pos+2:]
	args := []string{
		url,
		"-f", "bestaudio",
		"-x",
		"--audio-format", "mp3",
		"-o", "tmp/%(id)s.%(ext)s",
		"--cookies-from-browser", "brave",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"--no-playlist",
		"--sleep-interval", "5",
		"--max-sleep-interval", "15",
		"--sleep-requests", "1",
		"--extractor-retries", "3",
		"--retries", "3",
	}

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return "", err
	}

	_, err := dl.Run(ctx, args...)
	if err != nil {
		logger.Errorf("error running ytdlp %s", err)
		return "", err
	}
	return id, nil
}
func Search(ctx context.Context, query string) (ytdlp.ExtractedInfo, error) {
	var builder strings.Builder

	builder.WriteString("ytsearch1:")
	builder.WriteString(query)

	result := builder.String()

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return ytdlp.ExtractedInfo{}, err
	}

	output, err := dl.Run(ctx, result, "--print-json", "--skip-download", "--no-playlist")

	if err != nil {
		logger.Errorf("error running ytdlp %s", err)
		return ytdlp.ExtractedInfo{}, err
	}

	if len(output.Stdout) == 0 {
		logger.Errorf("yt-dlp returned no JSON (likely no results or rate-limited)")
		return ytdlp.ExtractedInfo{}, fmt.Errorf("Error")
	}

	var info ytdlp.ExtractedInfo
	if err := json.Unmarshal([]byte(output.Stdout), &info); err != nil {
		logger.Errorf("error unmarshaling %s", err)
		return ytdlp.ExtractedInfo{}, err
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
